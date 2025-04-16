package service

import (
	status_api "event_service/gen/status"
	"event_service/internal/models"
	"event_service/internal/repositories"
	"github.com/go-pg/pg/v10"
)

type StatusService struct {
	repo *repositories.StatusRepository
	db   *pg.DB
}

func NewStatusService(repo *repositories.StatusRepository, db *pg.DB) *StatusService {
	return &StatusService{
		repo: repo,
		db:   db,
	}
}

func SingleStatusConvert(model *models.Status) *status_api.StatusResponse {
	return &status_api.StatusResponse{
		Title: model.Title,
	}
}

func MultipleStatusConvert(models []*models.Status) []*status_api.StatusResponse {
	responses := make([]*status_api.StatusResponse, 0)

	for _, model := range models {
		responses = append(responses, SingleStatusConvert(model))
	}

	return responses
}

func (s *StatusService) GetAllStatuses() (_ []*status_api.StatusResponse, err error) {
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

	statusModels, err := s.repo.GetAllStatuses(tx)
	if err != nil {
		return nil, err
	}

	return MultipleStatusConvert(statusModels), nil
}

func (s *StatusService) GetStatusById(statusId int) (_ *status_api.StatusResponse, err error) {
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

	statusModel, err := s.repo.GetStatusById(tx, statusId)
	if err != nil {
		return nil, err
	}

	return SingleStatusConvert(statusModel), nil
}

func (s *StatusService) CreateStatus(status status_api.Status) (_ *status_api.StatusResponse, err error) {
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

	model := &models.Status{
		Title: status.Title,
	}

	statusModel, err := s.repo.Create(tx, model)
	if err != nil {
		return nil, err
	}

	return SingleStatusConvert(statusModel), nil
}

func (s *StatusService) UpdateStatus(eventId int, newStatus status_api.StatusUpdate) (_ *status_api.StatusResponse, err error) {
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

	status, err := s.GetStatusById(eventId)
	if err != nil {
		return status, err
	}

	if newStatus.Title != nil {
		status.Title = *newStatus.Title
	}

	model := &models.Status{
		ID:    eventId,
		Title: status.Title,
	}

	statusModel, err := s.repo.UpdateStatus(tx, model)
	if err != nil {
		return nil, err
	}

	return SingleStatusConvert(statusModel), nil
}

func (s *StatusService) DeleteStatus(statusId int) (err error) {
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

	return s.repo.DeleteStatus(tx, statusId)
}
