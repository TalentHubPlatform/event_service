package service

import (
	"event_service/internal/models"
	"event_service/internal/repositories"
	"event_service/internal/schemas"
	"github.com/go-pg/pg/v10"
)

type TimelineService struct {
	repo               *repositories.TimelineRepository
	timelineStatusRepo *repositories.TimelineStatusRepository
	db                 *pg.DB
}

func NewTimelineService(repo *repositories.TimelineRepository, timelineStatusRepo *repositories.TimelineStatusRepository,
	db *pg.DB) *TimelineService {
	return &TimelineService{
		repo:               repo,
		timelineStatusRepo: timelineStatusRepo,
		db:                 db,
	}
}

func (s *TimelineService) GetAllTimelines(trackId int) (_ []*models.Timeline, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return s.repo.GetTimelinesByTrackID(tx, trackId)
}

func (s *TimelineService) GetAllTimelinesWithStatus(trackId int, Status string) (_ []*models.Timeline, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return s.repo.GetTimelinesByTrackIDWithStatus(tx, trackId, Status)
}

func (s *TimelineService) GetTimelineById(timelineId int) (_ *models.Timeline, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return s.repo.GetTimelineByID(tx, timelineId)
}

func (s *TimelineService) CreateTimeline(timeline *schemas.Timeline) (_ *models.Timeline, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	nextIndex, err := s.repo.GetMaxNumOfTimeline(tx, timeline.TrackID)
	nextIndex++

	model := &models.Timeline{
		Title:            timeline.Title,
		Description:      timeline.Description,
		Deadline:         timeline.Deadline,
		IsBlocking:       timeline.IsBlocking,
		Status:           timeline.Status,
		TrackID:          timeline.TrackID,
		TimelineStatusID: nextIndex,
	}

	return s.repo.Create(tx, model)
}

func (s *TimelineService) UpdateTimeline(timelineId int, newTimeline *schemas.TimelineUpdate) (_ *models.Timeline, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	timeline, err := s.GetTimelineById(timelineId)
	if err != nil {
		return nil, err
	}

	if newTimeline.Title != "" {
		timeline.Title = newTimeline.Title
	}

	if newTimeline.Description != "" {
		timeline.Description = newTimeline.Description
	}

	if !newTimeline.Deadline.IsZero() {
		timeline.Deadline = newTimeline.Deadline
	}

	timeline.IsBlocking = newTimeline.IsBlocking

	if newTimeline.TrackID != 0 {
		timeline.TrackID = newTimeline.TrackID
	}

	if newTimeline.TimelineStatusID != 0 {
		timeline.TimelineStatusID = newTimeline.TimelineStatusID
	}

	if newTimeline.Status != "" {
		timeline.Status = newTimeline.Status
	}

	return s.repo.UpdateTimeline(tx, timelineId, timeline)
}

func (s *TimelineService) DeleteTimeline(timelineId int) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return s.repo.DeleteTimeline(tx, timelineId)
}

func (s *TimelineService) GetAllTimelineStatuses() (_ []*models.TimelineStatus, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	return s.timelineStatusRepo.GetAllTimelineStatuses(tx)
}

func (s *TimelineService) CreateTimelineStatus(timelineStatus *schemas.TimelineStatus) (_ *models.TimelineStatus, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	model := models.TimelineStatus{
		CountNum: timelineStatus.CountNum,
	}

	return s.timelineStatusRepo.Create(tx, &model)
}
