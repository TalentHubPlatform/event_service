package service

import (
	timeline_api "event_service/gen/timeline"
	"event_service/internal/models"
	"event_service/internal/repositories"
	"github.com/go-pg/pg/v10"
	"strconv"
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

func SingleTimelineConvert(model *models.Timeline) *timeline_api.TimelineResponse {
	return &timeline_api.TimelineResponse{
		Deadline:         model.Deadline,
		Description:      model.Description,
		IsBlocking:       model.IsBlocking,
		Status:           model.Status,
		Title:            model.Title,
		TrackId:          model.TrackID,
		TimelineStatusId: model.TimelineStatusID,
	}
}

func MultipleTimelineConvert(models []*models.Timeline) []*timeline_api.TimelineResponse {
	responses := make([]*timeline_api.TimelineResponse, 0)

	for _, model := range models {
		responses = append(responses, SingleTimelineConvert(model))
	}

	return responses
}

func MultipleTimelineStatusConvert(statusesModels []*models.TimelineStatus) []*timeline_api.TimelineStatusResponse {
	responses := make([]*timeline_api.TimelineStatusResponse, 0)

	for _, model := range statusesModels {
		responses = append(responses, &timeline_api.TimelineStatusResponse{
			CountNum: model.CountNum,
		})
	}

	return responses
}

func (s *TimelineService) GetAllTimelines(trackId int) (_ []*timeline_api.TimelineResponse, err error) {
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

	timelineModels, err := s.repo.GetTimelinesByTrackID(tx, trackId)
	if err != nil {
		return nil, err
	}

	return MultipleTimelineConvert(timelineModels), nil
}

func (s *TimelineService) GetAllTimelinesWithStatus(trackId int, Status string) (_ []*timeline_api.TimelineResponse, err error) {
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

	timelineModels, err := s.repo.GetTimelinesByTrackIDWithStatus(tx, trackId, Status)
	if err != nil {
		return nil, err
	}

	return MultipleTimelineConvert(timelineModels), err
}

func (s *TimelineService) GetTimelineById(timelineId int) (_ *timeline_api.TimelineResponse, err error) {
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

	timelineModel, err := s.repo.GetTimelineByID(tx, timelineId)
	if err != nil {
		return nil, err
	}

	return SingleTimelineConvert(timelineModel), nil
}

func (s *TimelineService) CreateTimeline(timeline *timeline_api.Timeline) (_ *timeline_api.TimelineResponse, err error) {
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

	nextIndex, err := s.repo.GetMaxNumOfTimeline(tx, timeline.TrackId)
	nextIndex++

	model := &models.Timeline{
		Title:            timeline.Title,
		Description:      timeline.Description,
		Deadline:         timeline.Deadline,
		IsBlocking:       timeline.IsBlocking,
		Status:           timeline.Status,
		TrackID:          timeline.TrackId,
		TimelineStatusID: nextIndex,
	}

	timelineModel, err := s.repo.Create(tx, model)
	if err != nil {
		return nil, err
	}

	return SingleTimelineConvert(timelineModel), nil
}

func (s *TimelineService) UpdateTimeline(timelineId int, newTimeline *timeline_api.TimelineUpdate) (_ *timeline_api.TimelineResponse, err error) {
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

	if newTimeline.Title != nil {
		timeline.Title = *newTimeline.Title
	}

	if newTimeline.Description != nil {
		timeline.Description = *newTimeline.Description
	}

	if !newTimeline.Deadline.IsZero() {
		timeline.Deadline = *newTimeline.Deadline
	}

	converted, err := strconv.ParseBool(*newTimeline.IsBlocking)
	if err != nil {
		return nil, err
	}

	timeline.IsBlocking = converted

	if newTimeline.TrackId != nil {
		timeline.TrackId = *newTimeline.TrackId
	}

	if newTimeline.TimelineStatusId != nil {
		timeline.TimelineStatusId = *newTimeline.TimelineStatusId
	}

	if newTimeline.Status != nil {
		timeline.Status = *newTimeline.Status
	}

	model := &models.Timeline{
		Title:            timeline.Title,
		Description:      timeline.Description,
		Deadline:         timeline.Deadline,
		IsBlocking:       timeline.IsBlocking,
		Status:           timeline.Status,
		TrackID:          timeline.TrackId,
		TimelineStatusID: timeline.TimelineStatusId,
	}

	timelineModel, err := s.repo.UpdateTimeline(tx, timelineId, model)
	if err != nil {
		return nil, err
	}

	return SingleTimelineConvert(timelineModel), nil
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

func (s *TimelineService) GetAllTimelineStatuses() (_ []*timeline_api.TimelineStatusResponse, err error) {
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

	timelineModels, err := s.timelineStatusRepo.GetAllTimelineStatuses(tx)
	if err != nil {
		return nil, err
	}

	return MultipleTimelineStatusConvert(timelineModels), nil
}

func (s *TimelineService) CreateTimelineStatus(timelineStatus *timeline_api.TimelineStatusResponse) (_ *timeline_api.TimelineStatusResponse, err error) {
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

	timelineModel, err := s.timelineStatusRepo.Create(tx, &model)
	if err != nil {
		return nil, err
	}

	return &timeline_api.TimelineStatusResponse{CountNum: timelineModel.CountNum}, nil
}
