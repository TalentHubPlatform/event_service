package service

import (
	api "event_service/gen/date"
	"event_service/internal/models"
	"event_service/internal/repositories"
	"github.com/go-pg/pg/v10"
)

type DateService struct {
	repo *repositories.DateRepository
	db   *pg.DB
}

func NewDateService(repo *repositories.DateRepository, db *pg.DB) *DateService {
	return &DateService{
		repo: repo,
		db:   db,
	}
}

func SingleDateConvert(model *models.Date) *api.DateResponse {
	return &api.DateResponse{
		DateStart: model.DateStart,
		DateEnd:   model.DateEnd,
		Id:        &model.ID,
	}
}

func MultipleDateConvert(models []*models.Date) []*api.DateResponse {
	responses := make([]*api.DateResponse, 0)

	for _, model := range models {
		responses = append(responses, SingleDateConvert(model))
	}

	return responses
}

func (s *DateService) GetAllDates() ([]*api.DateResponse, error) {
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

	dateModels, err := s.repo.GetAllDates(tx)
	if err != nil {
		return nil, err
	}

	return MultipleDateConvert(dateModels), nil
}

func (s *DateService) GetDateByID(id int) (*api.DateResponse, error) {
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

	dateModel, err := s.repo.GetDateById(tx, id)
	if err != nil {
		return nil, err
	}

	return SingleDateConvert(dateModel), nil
}

func (s *DateService) CreateDate(date api.Date) (*api.DateResponse, error) {
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

	model := models.Date{
		DateStart: date.DateStart,
		DateEnd:   date.DateEnd,
	}

	dateModel, err := s.repo.Create(tx, &model)
	if err != nil {
		return nil, err
	}

	return SingleDateConvert(dateModel), nil
}

func (s *DateService) UpdateDate(id int, date api.DateUpdate) (*api.DateResponse, error) {
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

	_, err = s.repo.ChangeDateStart(tx, id, date.DateStart)
	if err != nil {
		return nil, err
	}

	model, err := s.repo.ChangeDateEnd(tx, id, date.DateEnd)
	if err != nil {
		return nil, err
	}

	return SingleDateConvert(model), nil
}

func (s *DateService) DeleteDate(id int) error {
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

	return s.repo.DeleteDate(tx, id)
}
