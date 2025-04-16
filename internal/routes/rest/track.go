package rest

import (
	"encoding/json"
	"event_service/internal/models"
	"event_service/internal/schemas"
	"event_service/internal/service"
	"event_service/pkg/http/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type TrackService interface {
	GetAllTracks() ([]*models.Track, error)
	GetTrackById(int) (*models.Track, error)
	CreateTrack(schemas.Track) (*models.Track, error)
	UpdateTrack(int, schemas.TrackUpdate) (*models.Track, error)
	DeleteTrack(int) error

	GetAllTrackLocations(int) ([]*models.Location, error)
	AddLocationToTrack(*schemas.LocationTrack) (*models.LocationTrack, error)
	RemoveLocationFromTrack(*schemas.LocationTrack) error

	GetRegisteredTeams(int) ([]*models.TrackTeam, error)
	GetCertainRegisteredTeam(int, int) (*models.TrackTeam, error)
	RegisterTeam(*schemas.TrackTeam) (*models.TrackTeam, error)
	UpdateRegisteredTeam(int, int, schemas.TrackTeamUpdate) (*models.TrackTeam, error)
	DeleteRegisteredTeam(int, int) error
}

func NewTrack(log *slog.Logger, service *service.TrackService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllTracksHandler(log, service))
		r.Post("/", createTrackHandler(log, service, validate))

		r.Route("/location", func(r chi.Router) {
			r.Post("/", addLocationToTrackHandler(log, service))

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", getAllTrackLocationsHandler(log, service))
				r.Delete("/", removeLocationFromTrackHandler(log, service))
			})
		})

		r.Route("/team", func(r chi.Router) {
			r.Get("/", getRegisteredTeamsHandler(log, service))
			r.Post("/", registerTeamHandler(log, service, validate))
		})

		//r.Route("/judge", func(r chi.Router) {
		//	r.Get("/", getWinnerRecordsHandler(log, service))
		//	r.Post("/", createWinnerRecordHandler(log, service, validate))
		//	r.Put("/", changeWinnerRecordHandler(log, service, validate))
		//	r.Delete("/", deleteWinnerRecordHandler(log, service))
		//})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getTracksByIDHandler(log, service))
			r.Put("/", updateTrackHandler(log, service, validate))
			r.Delete("/", deleteTrackHandler(log, service))

			r.Route("/team/{teamId}", func(r chi.Router) {
				r.Put("/", updateRegisteredTeamHandler(log, service, validate))
				r.Delete("/", deleteRegisteredTeamHandler(log, service))
			})
		})
	})

	return r
}

func getAllTracksHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.getAll"

		log := log.With(
			slog.With("op", op),
			slog.With("request_id", middleware.GetReqID(r.Context())),
		)

		tracks, err := service.GetAllTracks()
		if err != nil {
			log.Error("error getting all tracks:", err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tracks); err != nil {
			log.Error("error encoding tracks:", err)

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("tracks successfully fetched")
	}
}

func getTracksByIDHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.getByID"

		log := log.With(
			slog.String("op", op),
			slog.String("request_it", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		track, err := service.GetTrackById(trackId)
		if err != nil {
			log.Error("Failed to get track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(track); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Track fetched successfully")
	}
}

func createTrackHandler(log *slog.Logger, service TrackService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var track schemas.Track
		if err := DecodeAndValidate(r, &track, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateTrack(track)
		if err != nil {
			log.Error("Failed to create track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Track created successfully")
	}
}

func updateTrackHandler(log *slog.Logger, service TrackService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.update"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))

		var track schemas.TrackUpdate
		if err := DecodeAndValidate(r, &track, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpdateTrack(trackId, track)
		if err != nil {
			log.Error("Failed to update track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Track updated successfully")
	}
}

func deleteTrackHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		err = service.DeleteTrack(trackId)

		if err != nil {
			log.Error("Failed to delete track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Track deleted successfully")
	}
}

func getAllTrackLocationsHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.locationGet"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		locations, err := service.GetAllTrackLocations(trackId)

		if err != nil {
			log.Error("Failed to get locations by track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(locations); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Locations by track fetched successfully")
	}
}

func addLocationToTrackHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.addLocation"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		headersList := map[string]string{
			"TrackId":    "int",
			"LocationId": "int",
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Failed to validate headers:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newLocation, err := service.AddLocationToTrack(&schemas.LocationTrack{
			LocationId: convertedHeaders["LocationId"].(int),
			TrackId:    convertedHeaders["TrackId"].(int),
		})

		if err != nil {
			log.Error("Failed to add location to track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newLocation); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Location added to track successfully")
	}
}

func removeLocationFromTrackHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.removeLocation"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		queryParams := r.URL.Query()

		if queryParams.Get("location_id") == "" {
			log.Error("Missing 'location_id' in query")

			http.Error(w, "Missing 'location_id' in query", http.StatusBadRequest)
			return
		}

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		locationId, err := strconv.Atoi(queryParams.Get("location_id"))
		if err != nil {
			log.Error("Invalid format of location_id:", err.Error())

			http.Error(w, fmt.Sprintf("Invalid format of location_id query, expected number, got %s", queryParams.Get("status_id")), http.StatusBadRequest)
			return
		}

		err = service.RemoveLocationFromTrack(&schemas.LocationTrack{
			TrackId:    trackId,
			LocationId: locationId,
		})

		if err != nil {
			log.Error("Failed to remove location from track:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Location deleted from track successfully")
	}
}

func getRegisteredTeamsHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.getRegisteredTeams"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		headersList := map[string]string{
			"TrackId": "int",
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Failed to validate headers: ", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		trackId := convertedHeaders["TrackId"].(int)

		registeredTeams, err := service.GetRegisteredTeams(trackId)
		if err != nil {
			log.Error("Failed to fetch registered teams:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(registeredTeams); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("RegisteredTeams successfully fetched")
	}
}

func registerTeamHandler(log *slog.Logger, service TrackService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.registerTeam"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var trackTeam schemas.TrackTeam
		if err := DecodeAndValidate(r, &trackTeam, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.RegisterTeam(&trackTeam)
		if err != nil {
			log.Error("Failed to create TrackTeam:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TrackTeam created successfully")
	}
}

func updateRegisteredTeamHandler(log *slog.Logger, service TrackService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.updateRegisteredTeam"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		teamId, err := strconv.Atoi(chi.URLParam(r, "teamId"))

		var trackTeam schemas.TrackTeamUpdate
		if err := DecodeAndValidate(r, &trackTeam, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpdateRegisteredTeam(trackId, teamId, trackTeam)
		if err != nil {
			log.Error("Failed to update registered team:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Registered team updated successfully")
	}
}

func deleteRegisteredTeamHandler(log *slog.Logger, service TrackService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Track.deleteRegisteredTeam"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "id"))
		teamId, err := strconv.Atoi(chi.URLParam(r, "teamId"))

		err = service.DeleteRegisteredTeam(trackId, teamId)

		if err != nil {
			log.Error("Failed to delete registered team:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Info("Registered team deleted successfully")
	}
}
