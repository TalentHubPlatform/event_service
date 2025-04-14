package service

import (
	locationapi "event_service/gen/location"
	"event_service/internal/models"
	"event_service/internal/repositories"
	"github.com/go-pg/pg/v10"
)

type LocationService struct {
	repo *repositories.LocationRepository
	db   *pg.DB
}

func SingleLocationConvert(model *models.Location) *locationapi.LocationResponse {
	return &locationapi.LocationResponse{
		Title: model.Title,
	}
}

func MultipleLocationConvert(models []*models.Location) []*locationapi.LocationResponse {
	responses := make([]*locationapi.LocationResponse, 0)

	for _, model := range models {
		responses = append(responses, SingleLocationConvert(model))
	}

	return responses
}

func NewLocationService(repo *repositories.LocationRepository, db *pg.DB) *LocationService {
	return &LocationService{
		repo: repo,
		db:   db,
	}
}

func (s *LocationService) GetAllLocations() (_ []*locationapi.LocationResponse, err error) {
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

	locationModels, err := s.repo.GetAllLocations(tx)
	if err != nil {
		return nil, err
	}

	return MultipleLocationConvert(locationModels), nil
}

func (s *LocationService) GetLocationById(locationId int) (_ *locationapi.LocationResponse, err error) {
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

	locationModel, err := s.repo.GetLocationById(tx, locationId)
	if err != nil {
		return nil, err
	}

	return SingleLocationConvert(locationModel), nil
}

func (s *LocationService) CreateLocation(location locationapi.Location) (_ *locationapi.LocationResponse, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	model := &models.Location{
		Title: location.Title,
	}

	locationModel, err := s.repo.Create(tx, model)
	if err != nil {
		return nil, err
	}

	return SingleLocationConvert(locationModel), nil
}

func (s *LocationService) UpdateLocation(locationId int, newLocation locationapi.LocationUpdate) (_ *locationapi.LocationResponse, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	location, err := s.GetLocationById(locationId)
	if err != nil {
		return nil, err
	}

	if newLocation.Title != nil {
		location.Title = *newLocation.Title
	}

	model := &models.Location{
		Title: location.Title,
	}

	locationModel, err := s.repo.Update(tx, locationId, model)
	if err != nil {
		return nil, err
	}

	return SingleLocationConvert(locationModel), nil
}

func (s *LocationService) DeleteLocation(locationId int) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		_ = tx.Commit()
	}()

	return s.repo.DeleteLocation(tx, locationId)
}
