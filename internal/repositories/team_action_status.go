package repositories

import (
	"event_service/internal/models"
	"github.com/go-pg/pg/v10"
	"sort"
)

type TeamActionStatusRepository struct {
	DB *pg.DB
}

type AggregateResult struct {
	TeamId     int
	TotalValue int
}

func NewTeamActionStatusRepository(db *pg.DB) *TeamActionStatusRepository {
	return &TeamActionStatusRepository{DB: db}
}

func (r *TeamActionStatusRepository) Create(tx *pg.Tx, teamActionStatus *models.TeamActionStatus) (*models.TeamActionStatus, error) {
	_, err := tx.Model(teamActionStatus).Insert()
	return teamActionStatus, err
}

func (r *TeamActionStatusRepository) GetTeamActionStatusByTeamID(tx *pg.Tx, teamID int) ([]*models.TeamActionStatus, error) {
	teamActionStatuses := make([]*models.TeamActionStatus, 0)
	err := tx.Model(&teamActionStatuses).Where("track_team_id = ?", teamID).Select()
	return teamActionStatuses, err
}

func (r *TeamActionStatusRepository) GetTeamActionStatusByTimelineID(tx *pg.Tx, timelineID int) ([]*models.TeamActionStatus, error) {
	teamActionStatuses := make([]*models.TeamActionStatus, 0)
	err := tx.Model(&teamActionStatuses).Where("timeline_id = ?", timelineID).Select()
	return teamActionStatuses, err
}

func (r *TeamActionStatusRepository) GetTeamActionStatusByTeamIDAndTimelineID(tx *pg.Tx, teamID, timelineID int) (*models.TeamActionStatus, error) {
	teamActionStatus := new(models.TeamActionStatus)
	err := tx.Model(teamActionStatus).Where("track_team_id = ?", teamID).Where("timeline_id = ?", timelineID).Select()
	return teamActionStatus, err
}

func (r *TeamActionStatusRepository) UpdateTeamActionStatus(tx *pg.Tx, teamID int, timelineID int, newTeamActionStatus *models.TeamActionStatus) (*models.TeamActionStatus, error) {
	teamActionStatus := new(models.TeamActionStatus)
	_, err := tx.Model(teamActionStatus).Set("result_value = ?, resolution_link = ?, completed_at = ?, notes = ?", newTeamActionStatus.ResultValue,
		newTeamActionStatus.ResolutionLink, newTeamActionStatus.CompletedAt, newTeamActionStatus.Notes).Where("timeline_id = ? AND track_team_id = ?", timelineID, teamID).Returning("*").Update()
	return teamActionStatus, err
}

func (r *TeamActionStatusRepository) DeleteTeamActionStatus(tx *pg.Tx, teamID int, timelineID int) error {
	teamActionStatus := new(models.TeamActionStatus)
	_, err := tx.Model(teamActionStatus).Where("track_team_id = ?", teamID).Where("timeline_id = ?", timelineID).Delete()
	return err
}

func (r *TeamActionStatusRepository) AggregateResults(tx *pg.Tx, trackId int, limit int, offset int) ([]*AggregateResult, error) {
	var results []*AggregateResult

	query := `
        SELECT 
            tas.track_team_id AS team_id,
            SUM(tas.result_value) AS total_value
        FROM 
            team_action_status tas
        JOIN 
        	timeline t
        ON
        	t.id = tas.timeline_id
    	WHERE
            t.track_id = ?
        GROUP BY 
            tas.track_team_id
        LIMIT ? OFFSET ?
    `

	_, err := tx.Query(&results, query, trackId, limit, offset)
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalValue > results[j].TotalValue
	})

	return results, nil
}
