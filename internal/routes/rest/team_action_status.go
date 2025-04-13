package rest

import (
	"encoding/json"
	"event_service/internal/models"
	"event_service/internal/schemas"
	"event_service/internal/service"
	"event_service/pkg/http/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type TeamActionStatusService interface {
	GetTeamActionStatusByTeamId(int) ([]*models.TeamActionStatus, error)
	GetTeamActionStatusByTimelineId(int) ([]*models.TeamActionStatus, error)
	GetTeamActionStatus(int, int) (*models.TeamActionStatus, error)
	CreateTeamActionStatus(*schemas.TeamActionStatus) (*models.TeamActionStatus, error)
	UpdateTeamActionStatus(int, int, *schemas.TeamActionStatusUpdate) (*models.TeamActionStatus, error)
	DeleteTeamActionStatus(int, int) error
}

func NewTeamActionStatus(log *slog.Logger, service *service.TeamActionStatusService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	r.Route("/", func(r chi.Router) {
		r.Get("/", getTeamActionStatusHandler(log, service))
		r.Post("/", createTeamActionStatusHandler(log, service, validate))

		r.Route("/{timelineId}/{teamId}", func(r chi.Router) {
			r.Get("/", getTeamActionStatusByIdHandler(log, service))
			r.Put("/", updateTeamActionStatusHandler(log, service, validate))
			r.Delete("/", deleteTeamActionStatusHandler(log, service))
		})
	})

	return r
}

func getTeamActionStatusHandler(log *slog.Logger, service TeamActionStatusService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TeamActionStatus.get"

		headersList := map[string]string{}
		if r.Header.Get("TeamId") != "" {
			headersList["TeamId"] = "int"
		} else {
			headersList["TimelineId"] = "int"
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Failed to validate headers:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result []*models.TeamActionStatus
		if _, has := headersList["TeamId"]; has {
			result, err = service.GetTeamActionStatusByTeamId(convertedHeaders["TeamId"].(int))
		} else {
			result, err = service.GetTeamActionStatusByTimelineId(convertedHeaders["TimelineId"].(int))
		}

		if err != nil {
			log.Error("Failed to get TeamActionStatuses:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TeamActionStatus successfully fetched")
	}
}

func createTeamActionStatusHandler(log *slog.Logger, service TeamActionStatusService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TeamActionStatus.create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var teamActionStatus schemas.TeamActionStatus
		if err := decodeAndValidate(r, &teamActionStatus, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateTeamActionStatus(&teamActionStatus)
		if err != nil {
			log.Error("Failed to create teamActionStatus:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TeamActionStatus created successfully")
	}
}

func getTeamActionStatusByIdHandler(log *slog.Logger, service TeamActionStatusService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TeamActionStatus.getById"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, _ := strconv.Atoi(chi.URLParam(r, "timelineId"))
		teamId, _ := strconv.Atoi(chi.URLParam(r, "teamId"))

		teamActionStatus, err := service.GetTeamActionStatus(timelineId, teamId)
		if err != nil {
			log.Error("Failed to get teamActionStatus by id:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(teamActionStatus); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TeamActionStatus successfully created")
	}
}

func updateTeamActionStatusHandler(log *slog.Logger, service TeamActionStatusService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TeamActionStatus.getById"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, _ := strconv.Atoi(chi.URLParam(r, "timelineId"))
		teamId, _ := strconv.Atoi(chi.URLParam(r, "teamId"))

		var teamActionStatus schemas.TeamActionStatusUpdate
		if err := decodeAndValidate(r, &teamActionStatus, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpdateTeamActionStatus(timelineId, teamId, &teamActionStatus)
		if err != nil {
			log.Error("Failed to update TeamActionStatus:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TeamActionStatus updated successfully")
	}
}

func deleteTeamActionStatusHandler(log *slog.Logger, service TeamActionStatusService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TeamActionStatus.delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, _ := strconv.Atoi(chi.URLParam(r, "timelineId"))
		teamId, _ := strconv.Atoi(chi.URLParam(r, "teamId"))

		err := service.DeleteTeamActionStatus(timelineId, teamId)
		if err != nil {
			log.Error("Failed to delete TeamActionStatus:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
