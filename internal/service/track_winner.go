package service

import (
	"event_service/internal/models"
	"event_service/internal/repositories"
	"event_service/internal/schemas"
	"fmt"
	"github.com/go-pg/pg/v10"
)

type TrackWinnerService struct {
	repo *repositories.TrackWinnerRepository

	teamActionStatusRepo *repositories.TeamActionStatusRepository
	trackRepo            *repositories.TrackRepository
	timelineRepo         *repositories.TimelineRepository

	db *pg.DB
}

func NewTrackWinnerService(repo *repositories.TrackWinnerRepository, trackRepo *repositories.TrackRepository,
	timelineRepo *repositories.TimelineRepository, db *pg.DB) *TrackWinnerService {
	return &TrackWinnerService{
		repo:         repo,
		trackRepo:    trackRepo,
		timelineRepo: timelineRepo,
		db:           db,
	}
}

func (s *TrackWinnerService) GetWinnersOfTrack(trackId int) (_ []*models.TrackWinner, err error) {
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

	return s.repo.GetAllWinnersByTrackID(tx, trackId)
}

func (s *TrackWinnerService) GetWinnerById(trackId int, trackTeamId int) (_ *models.TrackWinner, err error) {
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

	return s.repo.GetTrackWinnerByTrackIDAndTeamID(tx, trackId, trackTeamId)
}

func (s *TrackWinnerService) CreateWinnerOfTrack(trackWinner *schemas.TrackWinner) (_ *models.TrackWinner, err error) {
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

	track, err := s.trackRepo.GetTrackByID(tx, trackWinner.TrackID)
	if err != nil {
		return nil, err
	}

	if track.IsScoreBased {
		return nil, fmt.Errorf("manually creation of winners of track forbidden")
	}

	model := &models.TrackWinner{
		TrackID:     trackWinner.TrackID,
		TrackTeamID: trackWinner.TrackTeamID,
		Place:       trackWinner.Place,
		IsAwardee:   trackWinner.IsAwardee,
	}

	return s.repo.Create(tx, model)
}

func (s *TrackWinnerService) CalculateRating(trackId int, limit int, offset int) ([]*repositories.AggregateResult, error) {
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

	return s.teamActionStatusRepo.AggregateResults(tx, trackId, limit, offset)
}

func (s *TrackWinnerService) SetResultsOfTrack(trackId int, threshold int, limit int) ([]*models.TrackWinner, error) {
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

	track, err := s.trackRepo.GetTrackByID(tx, trackId)
	if err != nil {
		return nil, err
	}

	if !track.IsScoreBased || track.Status != "completed" {
		return nil, fmt.Errorf("results can be generated only for completed tracks based on score")
	}

	aggregated, err := s.CalculateRating(trackId, limit, 0)
	if err != nil {
		return nil, err
	}

	result := make([]*models.TrackWinner, 0)

	for idx, resultItem := range aggregated {
		model := &models.TrackWinner{
			TrackID:     trackId,
			TrackTeamID: resultItem.TeamId,
			Place:       idx,
			IsAwardee:   resultItem.TotalValue >= threshold,
		}

		response, err := s.repo.Create(tx, model)
		if err != nil {
			return nil, err
		}

		result = append(result, response)
	}

	return result, nil
}
