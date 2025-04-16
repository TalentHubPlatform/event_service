package repositories

import (
	"event_service/internal/models"
	"github.com/go-pg/pg/v10"
)

type TimelineRepository struct {
	DB *pg.DB
}

func NewTimelineRepository(db *pg.DB) *TimelineRepository {
	return &TimelineRepository{DB: db}
}

func (r *TimelineRepository) Create(tx *pg.Tx, timeline *models.Timeline) (*models.Timeline, error) {
	_, err := tx.Model(timeline).Insert()
	return timeline, err
}

func (r *TimelineRepository) GetTimelinesByTrackID(tx *pg.Tx, trackID int) ([]*models.Timeline, error) {
	timelines := make([]*models.Timeline, 0)

	err := tx.Model(&timelines).
		Where("track_id = ?", trackID).
		Join("JOIN timeline_status as ts ON ts.id = timeline.timeline_status_id").
		Order("ts.count_num").
		Select()

	return timelines, err
}

func (r *TimelineRepository) GetMaxNumOfTimeline(tx *pg.Tx, trackID int) (int, error) {
	maxCountNum := 0

	err := tx.Model((*models.Timeline)(nil)).
		ColumnExpr("MAX(ts.count_num)").
		Join("JOIN timeline_status AS ts ON ts.id = timeline.timeline_status_id").
		Where("timeline.track_id = ?", trackID).
		Select(&maxCountNum)

	return maxCountNum, err
}

func (r *TimelineRepository) GetTimelinesByTrackIDWithStatus(tx *pg.Tx, trackID int, Status string) ([]*models.Timeline, error) {
	timelines := make([]*models.Timeline, 0)

	err := tx.Model(&timelines).
		Where("track_id = ? AND status = ?", trackID, Status).
		Join("JOIN timeline_status as ts ON ts.id = timeline.timeline_status_id").
		Order("ts.count_num").
		Select()

	return timelines, err
}

func (r *TimelineRepository) GetTimelineByID(tx *pg.Tx, timelineID int) (*models.Timeline, error) {
	timeline := new(models.Timeline)
	err := tx.Model(timeline).Where("id = ?", timelineID).Select()
	return timeline, err
}

func (r *TimelineRepository) UpdateTimeline(tx *pg.Tx, timelineId int, newTimeline *models.Timeline) (*models.Timeline, error) {
	timeline := new(models.Timeline)
	_, err := tx.Model(timeline).Set("title = ?, description = ?, deadline = ?, is_blocking = ?, timeline_status_id", newTimeline.Title,
		newTimeline.Description, newTimeline.Deadline, newTimeline.IsBlocking).Where("id = ?", timelineId).Returning("*").Update()
	return timeline, err
}

func (r *TimelineRepository) DeleteTimeline(tx *pg.Tx, timelineID int) error {
	_, err := tx.Model(&models.Timeline{}).Where("id = ?", timelineID).Delete()
	return err
}

func (r *TimelineRepository) GetMaxValue(tx *pg.Tx, trackId int) (int, error) {
	query := `
        SELECT SUM(CASE WHEN is_scoring THEN 100 ELSE 0 END) AS total_value
        FROM timeline
        WHERE track_id = ?
    `

	var result int

	_, err := tx.QueryOne(pg.Scan(&result), query, trackId)
	if err != nil {
		return 0, err
	}

	return result, nil
}
