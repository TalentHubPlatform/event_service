package repositories

import (
	"event_service/internal/models"
	"github.com/go-pg/pg/v10"
)

type TrackTeamRepository struct {
	DB *pg.DB
}

func NewTrackTeamRepository(db *pg.DB) *TrackTeamRepository {
	return &TrackTeamRepository{DB: db}
}

func (r *TrackTeamRepository) Create(tx *pg.Tx, trackTeam *models.TrackTeam) (*models.TrackTeam, error) {
	_, err := tx.Model(trackTeam).Insert()
	return trackTeam, err
}

func (r *TrackTeamRepository) GetTeamsByTrackID(tx *pg.Tx, trackID int) ([]*models.TrackTeam, error) {
	trackTeams := make([]*models.TrackTeam, 0)
	err := tx.Model(&trackTeams).Where("track_id = ?", trackID).Select()
	return trackTeams, err
}

func (r *TrackTeamRepository) GetTracksByTeamID(tx *pg.Tx, teamID int) ([]*models.TrackTeam, error) {
	trackTeams := make([]*models.TrackTeam, 0)
	err := tx.Model(&trackTeams).Where("team_id = ?", teamID).Select()
	return trackTeams, err
}

func (r *TrackTeamRepository) GetByTrackIDAndTeamID(tx *pg.Tx, trackID, teamID int) (*models.TrackTeam, error) {
	trackTeam := new(models.TrackTeam)
	err := tx.Model(trackTeam).Where("track_id = ?", trackID).Where("team_id = ?", teamID).Select()
	return trackTeam, err
}

func (r *TrackTeamRepository) UpdateTrackTeam(tx *pg.Tx, trackID, teamID int, trackTeam *models.TrackTeam) (*models.TrackTeam, error) {
	_, err := tx.Model(trackTeam).Set("is_active = ?", trackTeam.IsActive).Where("track_id = ?", trackID).Where("team_id = ?", teamID).Update()
	return trackTeam, err
}

func (r *TrackTeamRepository) DeleteTrackTeam(tx *pg.Tx, trackID, teamID int) error {
	_, err := tx.Model(&models.TrackTeam{}).Where("track_id = ?", trackID).Where("team_id = ?", teamID).Delete()
	return err
}
