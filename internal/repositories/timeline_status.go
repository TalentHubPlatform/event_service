package repositories

import (
	"event_service/internal/models"
	"github.com/go-pg/pg/v10"
)

type TimelineStatusRepository struct {
	DB *pg.DB
}

func NewTimelineStatusRepository(db *pg.DB) *TimelineStatusRepository {
	return &TimelineStatusRepository{DB: db}
}

func (r *TimelineStatusRepository) Create(tx *pg.Tx, timelineStatus *models.TimelineStatus) (*models.TimelineStatus, error) {
	_, err := tx.Model(timelineStatus).Insert()
	return timelineStatus, err
}

func (r *TimelineStatusRepository) GetAllTimelineStatuses(tx *pg.Tx) ([]*models.TimelineStatus, error) {
	timelineStatuses := make([]*models.TimelineStatus, 0)
	err := tx.Model(&timelineStatuses).Select()
	return timelineStatuses, err
}
