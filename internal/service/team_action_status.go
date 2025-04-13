package service

import (
	"event_service/internal/models"
	"event_service/internal/repositories"
	"event_service/internal/schemas"
	"github.com/go-pg/pg/v10"
	"strconv"
)

type TeamActionStatusService struct {
	db   *pg.DB
	repo *repositories.TeamActionStatusRepository
}

func NewTeamActionStatusService(repo *repositories.TeamActionStatusRepository, db *pg.DB) *TeamActionStatusService {
	return &TeamActionStatusService{
		repo: repo,
		db:   db,
	}
}

func (s *TeamActionStatusService) GetTeamActionStatusByTeamId(teamId int) (_ []*models.TeamActionStatus, err error) {
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

	return s.repo.GetTeamActionStatusByTeamID(tx, teamId)
}

func (s *TeamActionStatusService) GetTeamActionStatusByTimelineId(timelineId int) (_ []*models.TeamActionStatus, err error) {
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

	return s.repo.GetTeamActionStatusByTeamID(tx, timelineId)
}

func (s *TeamActionStatusService) GetTeamActionStatus(timelineId int, teamId int) (_ *models.TeamActionStatus, err error) {
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

	return s.repo.GetTeamActionStatusByTeamIDAndTimelineID(tx, teamId, timelineId)
}

func (s *TeamActionStatusService) CreateTeamActionStatus(teamActionStatus *schemas.TeamActionStatus) (_ *models.TeamActionStatus, err error) {
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

	model := &models.TeamActionStatus{
		TrackTeamID:    teamActionStatus.TrackTeamID,
		TimelineID:     teamActionStatus.TimelineID,
		ResolutionLink: teamActionStatus.ResolutionLink,
		ResultValue:    teamActionStatus.ResultValue,
		CompletedAt:    teamActionStatus.CompletedAt,
		Notes:          teamActionStatus.Notes,
	}

	return s.repo.Create(tx, model)
}

func (s *TeamActionStatusService) UpdateTeamActionStatus(timelineId int, teamId int, newTeamActionStatus *schemas.TeamActionStatusUpdate) (_ *models.TeamActionStatus, err error) {
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

	teamActionStatus, err := s.GetTeamActionStatus(timelineId, teamId)
	if err != nil {
		return nil, err
	}

	if newTeamActionStatus.ResolutionLink != "" {
		teamActionStatus.ResolutionLink = newTeamActionStatus.ResolutionLink
	}

	if newTeamActionStatus.Notes != "" {
		teamActionStatus.Notes = newTeamActionStatus.Notes
	}

	if !newTeamActionStatus.CompletedAt.IsZero() {
		teamActionStatus.CompletedAt = newTeamActionStatus.CompletedAt
	}

	if newTeamActionStatus.ResultValue != "" {
		converted, err := strconv.Atoi(newTeamActionStatus.ResultValue)
		if err != nil {
			return nil, err
		}

		teamActionStatus.ResultValue = converted
	}

	return s.repo.UpdateTeamActionStatus(tx, teamId, timelineId, teamActionStatus)
}

func (s *TeamActionStatusService) DeleteTeamActionStatus(timelineId int, teamId int) (err error) {
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

	return s.repo.DeleteTeamActionStatus(tx, teamId, timelineId)
}
