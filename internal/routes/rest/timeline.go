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

type TimelineService interface {
	GetAllTimelines(trackId int) ([]*models.Timeline, error)
	GetAllTimelinesWithStatus(int, string) ([]*models.Timeline, error)
	GetTimelineById(int) (*models.Timeline, error)
	CreateTimeline(*schemas.Timeline) (*models.Timeline, error)
	UpdateTimeline(int, *schemas.TimelineUpdate) (*models.Timeline, error)
	DeleteTimeline(int) error

	GetAllTimelineStatuses() ([]*models.TimelineStatus, error)
	CreateTimelineStatus(*schemas.TimelineStatus) (*models.TimelineStatus, error)
}

func NewTimeline(log *slog.Logger, service *service.TimelineService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	r.Route("/", func(r chi.Router) {
		r.Post("/", createTimelineHandler(log, service, validate))
		r.Get("/", GetTimelinesHandler(log, service))

		r.Route("/status", func(r chi.Router) {
			r.Get("/", GetTimelineStatusesHandler(log, service))
			r.Post("/", CreateTimelineStatusHandler(log, service, validate))
		})

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetTimelineByIdHandler(log, service))
			r.Put("/", UpdateTimelineHandler(log, service, validate))
			r.Delete("/", DeleteTimelineHandler(log, service))
		})
	})

	return r
}

func GetTimelinesHandler(log *slog.Logger, service TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Timeline.get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		headersList := map[string]string{
			"TrackId": "int",
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Failed to validate headers:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		queryParams := r.URL.Query()
		var result []*models.Timeline

		trackId := convertedHeaders["TrackId"].(int)
		status := queryParams.Get("status")

		if status == "" {
			result, err = service.GetAllTimelines(trackId)
		} else {
			result, err = service.GetAllTimelinesWithStatus(trackId, status)
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Timeline created fetched")
	}
}

func createTimelineHandler(log *slog.Logger, service TimelineService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Timeline.create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var timeline schemas.Timeline
		if err := decodeAndValidate(r, &timeline, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateTimeline(&timeline)
		if err != nil {
			log.Error("Failed to create timeline:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Timeline created successfully")
	}
}

func GetTimelineByIdHandler(log *slog.Logger, service TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Timeline.getById"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, err := strconv.Atoi(chi.URLParam(r, "id"))
		timeline, err := service.GetTimelineById(timelineId)

		if err != nil {
			log.Error("Failed to get timeline by id:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(timeline); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Timeline successfully fetched by id")
	}
}

func UpdateTimelineHandler(log *slog.Logger, service TimelineService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Timeline.getById"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, err := strconv.Atoi(chi.URLParam(r, "id"))

		var timeline schemas.TimelineUpdate
		if err := decodeAndValidate(r, &timeline, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpdateTimeline(timelineId, &timeline)
		if err != nil {
			log.Error("Failed to update timeline:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Timeline updated successfully")
	}
}

func DeleteTimelineHandler(log *slog.Logger, service TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.Timeline.delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		timelineId, err := strconv.Atoi(chi.URLParam(r, "id"))
		err = service.DeleteTimeline(timelineId)

		if err != nil {
			log.Error("Failed to delete timeline by id:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetTimelineStatusesHandler(log *slog.Logger, service TimelineService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TimelineStatuses.get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		resp, err := service.GetAllTimelineStatuses()
		if err != nil {
			log.Error("Failed to get TimelineStatuses", err.Error())

			http.Error(w, "Failed to get TimelineStatuses", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TimelineStatuses successfully fetched")
	}
}

func CreateTimelineStatusHandler(log *slog.Logger, service TimelineService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TimelineStatus.create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var timelineStatus schemas.TimelineStatus
		if err := decodeAndValidate(r, &timelineStatus, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateTimelineStatus(&timelineStatus)
		if err != nil {
			log.Error("Failed to create timelineStatus:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TimelineStatus created successfully")
	}
}
