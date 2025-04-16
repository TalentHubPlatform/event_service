package rest

import (
	status_api "event_service/gen/status"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type StatusService interface {
	GetAllStatuses() ([]*status_api.StatusResponse, error)
	GetStatusById(id int) (*status_api.StatusResponse, error)
	CreateStatus(date status_api.Status) (*status_api.StatusResponse, error)
	UpdateStatus(id int, date status_api.StatusUpdate) (*status_api.StatusResponse, error)
	DeleteStatus(id int) error
}

type StatusHandler struct {
	log       *slog.Logger
	service   StatusService
	validator *validator.Validate
}

func NewStatusHandler(log *slog.Logger, service StatusService, validator *validator.Validate) *StatusHandler {
	return &StatusHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *StatusHandler) GetStatus(ctx echo.Context) error {
	const op = "rest.Status.getAll"

	log := h.log.With(
		slog.String("op", op),
	)

	statuses, err := h.service.GetAllStatuses()
	if err != nil {
		log.Error("Failed to get statuses:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get statuses",
		})
	}

	log.Info("Statuses fetched successfully")

	return ctx.JSON(http.StatusOK, statuses)
}

func (h *StatusHandler) GetStatusId(ctx echo.Context, id status_api.Id) error {
	const op = "rest.Status.getByID"

	log := h.log.With(
		slog.String("op", op),
		slog.String("id", strconv.Itoa(id)),
	)

	status, err := h.service.GetStatusById(int(id))
	if err != nil {
		log.Error("Failed to get status:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get status",
		})
	}

	log.Info("Status fetched successfully")

	return ctx.JSON(http.StatusOK, status)
}

func (h *StatusHandler) PostStatus(ctx echo.Context) error {
	const op = "rest.Status.create"

	log := h.log.With(
		slog.String("op", op),
	)

	var status status_api.Status
	if err := decodeAndValidateEcho(ctx, &status, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.CreateStatus(status)
	if err != nil {
		log.Error("Failed to create status:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create status",
		})
	}

	log.Info("Status created successfully")

	return ctx.JSON(http.StatusCreated, resp)
}

func (h *StatusHandler) PutStatusId(ctx echo.Context, id status_api.Id) error {
	const op = "rest.Status.update"

	log := h.log.With(
		slog.String("op", op),
	)

	var status status_api.StatusUpdate
	if err := decodeAndValidateEcho(ctx, &status, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.UpdateStatus(int(id), status)
	if err != nil {
		log.Error("Failed to update status:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update status",
		})
	}

	log.Info("Status updated successfully")

	return ctx.JSON(http.StatusOK, resp)
}

func (h *StatusHandler) DeleteStatusId(ctx echo.Context, id status_api.Id) error {
	const op = "rest.Status.delete"

	log := h.log.With(
		slog.String("op", op),
	)

	err := h.service.DeleteStatus(int(id))
	if err != nil {
		log.Error("Failed to delete status:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete status",
		})
	}

	log.Info("Status deleted successfully")

	return ctx.NoContent(http.StatusOK)
}

func NewStatus(log *slog.Logger, service StatusService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	handler := NewStatusHandler(log, service, validate)

	r.Route("/", func(r chi.Router) {
		r.Get("/", HandlerAdapter(handler.GetStatus))
		r.Post("/", HandlerAdapter(handler.PostStatus))

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				log.Info(ctx.Param("Id"))
				return handler.GetStatusId(ctx, status_api.Id(id))
			}))

			r.Put("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				return handler.PutStatusId(ctx, status_api.Id(id))
			}))

			r.Delete("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				return handler.DeleteStatusId(ctx, status_api.Id(id))
			}))
		})
	})

	return r
}
