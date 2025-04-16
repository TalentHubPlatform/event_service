package rest

import (
	timeline_api "event_service/gen/timeline"
	"event_service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type TimelineService interface {
	GetAllTimelines(trackId int) ([]*timeline_api.TimelineResponse, error)
	GetAllTimelinesWithStatus(int, string) ([]*timeline_api.TimelineResponse, error)
	GetTimelineById(int) (*timeline_api.TimelineResponse, error)
	CreateTimeline(*timeline_api.Timeline) (*timeline_api.TimelineResponse, error)
	UpdateTimeline(int, *timeline_api.TimelineUpdate) (*timeline_api.TimelineResponse, error)
	DeleteTimeline(int) error

	GetAllTimelineStatuses() ([]*timeline_api.TimelineStatusResponse, error)
	CreateTimelineStatus(response *timeline_api.TimelineStatusResponse) (*timeline_api.TimelineStatusResponse, error)
}

type TimelineHandler struct {
	log       *slog.Logger
	service   TimelineService
	validator *validator.Validate
}

func NewTimelineHandler(log *slog.Logger, service TimelineService, validator *validator.Validate) *TimelineHandler {
	return &TimelineHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *TimelineHandler) GetTimeline(ctx echo.Context, params timeline_api.GetTimelineParams) error {
	const op = "rest.Timeline.get"

	log := h.log.With(
		slog.String("op", op),
	)

	if params.XTrackId == nil {
		log.Error("Missing required header: XTrackId")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required header: XTrackId",
		})
	}

	trackId := *params.XTrackId
	status := ""
	if params.Status != nil {
		status = *params.Status
	}

	var result []*timeline_api.TimelineResponse
	var err error

	if status == "" {
		result, err = h.service.GetAllTimelines(trackId)
	} else {
		result, err = h.service.GetAllTimelinesWithStatus(trackId, status)
	}

	if err != nil {
		log.Error("Failed to get timelines:", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get timelines",
		})
	}

	log.Info("Timelines fetched successfully")

	return ctx.JSON(http.StatusOK, result)
}

func (h *TimelineHandler) PostTimeline(ctx echo.Context) error {
	const op = "rest.Timeline.create"

	log := h.log.With(
		slog.String("op", op),
	)

	var timeline timeline_api.Timeline
	if err := decodeAndValidateEcho(ctx, &timeline, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.CreateTimeline(&timeline)
	if err != nil {
		log.Error("Failed to create timeline:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create timeline",
		})
	}

	log.Info("Timeline created successfully")

	return ctx.JSON(http.StatusCreated, resp)
}

func (h *TimelineHandler) GetTimelineId(ctx echo.Context, id timeline_api.Id) error {
	const op = "rest.Timeline.getById"

	log := h.log.With(
		slog.String("op", op),
	)

	timeline, err := h.service.GetTimelineById(int(id))
	if err != nil {
		log.Error("Failed to get timeline by id:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get timeline",
		})
	}

	log.Info("Timeline fetched successfully")

	return ctx.JSON(http.StatusOK, timeline)
}

func (h *TimelineHandler) DeleteTimelineId(ctx echo.Context, id timeline_api.Id) error {
	const op = "rest.Timeline.delete"

	log := h.log.With(
		slog.String("op", op),
	)

	err := h.service.DeleteTimeline(int(id))
	if err != nil {
		log.Error("Failed to delete timeline:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete timeline",
		})
	}

	log.Info("Timeline deleted successfully")

	return ctx.NoContent(http.StatusOK)
}

func (h *TimelineHandler) PutTimelineId(ctx echo.Context, id timeline_api.Id) error {
	const op = "rest.Timeline.update"

	log := h.log.With(
		slog.String("op", op),
	)

	var timeline timeline_api.TimelineUpdate
	if err := decodeAndValidateEcho(ctx, &timeline, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.UpdateTimeline(int(id), &timeline)
	if err != nil {
		log.Error("Failed to update timeline:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update timeline",
		})
	}

	log.Info("Timeline updated successfully")

	return ctx.JSON(http.StatusOK, resp)
}

func (h *TimelineHandler) GetTimelimeStatus(ctx echo.Context) error {
	const op = "rest.TimelineStatuses.get"

	log := h.log.With(
		slog.String("op", op),
	)

	resp, err := h.service.GetAllTimelineStatuses()
	if err != nil {
		log.Error("Failed to get TimelineStatuses:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get TimelineStatuses",
		})
	}

	log.Info("TimelineStatuses successfully fetched")

	return ctx.JSON(http.StatusOK, resp)
}

func (h *TimelineHandler) PostTimelimeStatus(ctx echo.Context) error {
	const op = "rest.TimelineStatus.create"

	log := h.log.With(
		slog.String("op", op),
	)

	var timelineStatus timeline_api.TimelineStatusResponse
	if err := decodeAndValidateEcho(ctx, &timelineStatus, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.CreateTimelineStatus(&timelineStatus)
	if err != nil {
		log.Error("Failed to create TimelineStatus:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create TimelineStatus",
		})
	}

	log.Info("TimelineStatus created successfully")

	return ctx.JSON(http.StatusCreated, resp)
}

func NewTimeline(log *slog.Logger, service *service.TimelineService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	handler := NewTimelineHandler(log, service, validate)

	r.Route("/", func(r chi.Router) {
		r.Get("/", HandlerAdapter(func(ctx echo.Context) error {
			queryParams := ctx.QueryParams()
			headers := ctx.Request().Header

			params := timeline_api.GetTimelineParams{}
			if trackId := headers.Get("XTrackId"); trackId != "" {
				id, err := strconv.Atoi(trackId)
				if err != nil {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid XTrackId"})
				}
				params.XTrackId = &id
			} else {
				return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required header: XTrackId"})
			}

			if status := queryParams.Get("status"); status != "" {
				params.Status = &status
			}

			return handler.GetTimeline(ctx, params)
		}))

		r.Post("/", HandlerAdapter(handler.PostTimeline))

		r.Route("/status", func(r chi.Router) {
			r.Get("/", HandlerAdapter(handler.GetTimelimeStatus))

			r.Post("/", HandlerAdapter(handler.PostTimelimeStatus))
		})

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", HandlerAdapter(func(ctx echo.Context) error {
				id, err := strconv.Atoi(ctx.Param("Id"))
				if err != nil {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
				}
				return handler.GetTimelineId(ctx, timeline_api.Id(id))
			}))

			r.Put("/", HandlerAdapter(func(ctx echo.Context) error {
				id, err := strconv.Atoi(ctx.Param("Id"))
				if err != nil {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
				}
				return handler.PutTimelineId(ctx, timeline_api.Id(id))
			}))

			r.Delete("/", HandlerAdapter(func(ctx echo.Context) error {
				id, err := strconv.Atoi(ctx.Param("Id"))
				if err != nil {
					return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
				}
				return handler.DeleteTimelineId(ctx, timeline_api.Id(id))
			}))
		})
	})

	return r
}
