package service

import (
	"event_service/internal/models"
	"event_service/internal/repositories"
	"event_service/internal/schemas"
	"fmt"
	"github.com/go-pg/pg/v10"
	"strconv"
)

type TrackService struct {
	repo *repositories.TrackRepository

	locationTrackRepo *repositories.LocationTrackRepository
	trackTeamRepo     *repositories.TrackTeamRepository

	db *pg.DB
}

func NewTrackService(repo *repositories.TrackRepository, db *pg.DB) *TrackService {
	return &TrackService{
		repo: repo,
		db:   db,
	}
}

func (s *TrackService) GetAllTracks() ([]*models.Track, error) {
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

	return s.repo.GetAllTracks(tx)
}

func (s *TrackService) GetTrackById(trackId int) (_ *models.Track, err error) {
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

	return s.repo.GetTrackByID(tx, trackId)
}

func (s *TrackService) GetTracksByEventId(eventId int) (_ []*models.Track, err error) {
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

	return s.repo.GetAllTracksByEventID(tx, eventId)
}

func (s *TrackService) StartTrack(trackId int) (_ *models.Track, err error) {
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

	return s.UpdateTrack(trackId, schemas.TrackUpdate{Status: "in_process"})
}

func (s *TrackService) EndTrack(trackId int) (_ *models.Track, err error) {
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

	return s.UpdateTrack(trackId, schemas.TrackUpdate{Status: "completed"})
}

func (s *TrackService) GetAllTracksToStart() (_ []*models.Track, err error) {
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

	return s.repo.GetAllTracksToStart(tx)
}

func (s *TrackService) GetAllTracksToEnd() (_ []*models.Track, err error) {
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

	return s.repo.GetAllTracksToEnd(tx)
}

func (s *TrackService) CreateTrack(track schemas.Track) (_ *models.Track, err error) {
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

	trackModel := &models.Track{
		Title:        track.Title,
		Description:  track.Description,
		IsScoreBased: track.IsScoreBased,
		EventID:      track.EventID,
		DateID:       track.DateID,
	}

	return s.repo.Create(tx, trackModel)
}

func (s *TrackService) UpdateTrack(trackId int, newTrack schemas.TrackUpdate) (*models.Track, error) {
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

	track, err := s.GetTrackById(trackId)
	if err != nil {
		return nil, err
	}

	if newTrack.Title != "" {
		track.Title = newTrack.Title
	}

	if newTrack.Description != "" {
		track.Description = newTrack.Description
	}

	if newTrack.IsScoreBased != "" {
		isScoreBased, err := strconv.ParseBool(newTrack.IsScoreBased)
		if err != nil {
			return nil, fmt.Errorf("invalid value for IsScoreBased: %v", err)
		}

		track.IsScoreBased = isScoreBased
	}

	if newTrack.EventID != 0 {
		track.EventID = newTrack.EventID
	}

	if newTrack.DateID != 0 {
		track.DateID = newTrack.DateID
	}

	return s.repo.UpdateTrack(tx, trackId, track)
}

func (s *TrackService) DeleteTrack(trackId int) error {
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

	return s.repo.DeleteTrack(tx, trackId)
}

func (s *TrackService) GetAllTrackLocations(trackId int) (_ []*models.Location, err error) {
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

	return s.locationTrackRepo.GetAllTracksLocations(tx, trackId)
}

func (s *TrackService) AddLocationToTrack(locationTrackSchema *schemas.LocationTrack) (_ *models.LocationTrack, err error) {
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

	statusEventModel := &models.LocationTrack{
		TrackId:    locationTrackSchema.TrackId,
		LocationId: locationTrackSchema.LocationId,
	}

	return s.locationTrackRepo.Create(tx, statusEventModel)
}

func (s *TrackService) RemoveLocationFromTrack(locationTrackSchema *schemas.LocationTrack) (err error) {
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

	return s.locationTrackRepo.DeleteTrackLocation(tx, locationTrackSchema.TrackId, locationTrackSchema.LocationId)
}

func (s *TrackService) GetRegisteredTeams(trackId int) (_ []*models.TrackTeam, err error) {
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

	return s.trackTeamRepo.GetTeamsByTrackID(tx, trackId)
}

func (s *TrackService) RegisterTeam(trackTeam *schemas.TrackTeam) (_ *models.TrackTeam, err error) {
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

	model := &models.TrackTeam{
		TrackID:  trackTeam.TrackID,
		TeamID:   trackTeam.TeamID,
		IsActive: trackTeam.IsActive,
	}

	return s.trackTeamRepo.Create(tx, model)
}

func (s *TrackService) GetCertainRegisteredTeam(trackId int, teamId int) (_ *models.TrackTeam, err error) {
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

	return s.trackTeamRepo.GetRoleByTrackIDAndTeamID(tx, trackId, teamId)
}

func (s *TrackService) UpdateRegisteredTeam(trackId int, teamId int, newTrackTeam schemas.TrackTeamUpdate) (_ *models.TrackTeam, err error) {
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

	trackTeam, err := s.GetCertainRegisteredTeam(trackId, teamId)
	if err != nil {
		return nil, err
	}

	if newTrackTeam.IsActive != "" {
		converted, err := strconv.ParseBool(newTrackTeam.IsActive)
		if err != nil {
			return nil, err
		}

		trackTeam.IsActive = converted
	}

	return s.trackTeamRepo.UpdateTrackTeam(tx, trackId, teamId, trackTeam)
}

func (s *TrackService) DeleteRegisteredTeam(trackId int, teamId int) error {
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

	return s.trackTeamRepo.DeleteTrackTeam(tx, trackId, teamId)
}
