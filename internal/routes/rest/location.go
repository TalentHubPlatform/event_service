package rest

import (
	locationapi "event_service/gen/location"
	"event_service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type LocationService interface {
	GetAllLocations() ([]*locationapi.LocationResponse, error)
	GetLocationById(locationId int) (*locationapi.LocationResponse, error)
	CreateLocation(date locationapi.Location) (*locationapi.LocationResponse, error)
	UpdateLocation(locationId int, date locationapi.LocationUpdate) (*locationapi.LocationResponse, error)
	DeleteLocation(locationId int) error
}

type LocationHandler struct {
	log       *slog.Logger
	service   LocationService
	validator *validator.Validate
}

func NewLocationHandler(log *slog.Logger, service LocationService, validator *validator.Validate) *LocationHandler {
	return &LocationHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *LocationHandler) GetLocation(ctx echo.Context) error {
	const op = "rest.Location.getAll"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", ctx.Get("request_id").(string)),
	)

	locations, err := h.service.GetAllLocations()
	if err != nil {
		log.Error("Failed to get locations:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get locations",
		})
	}

	log.Info("Locations fetched successfully")

	return ctx.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) GetLocationId(ctx echo.Context, id locationapi.Id) error {
	const op = "rest.Location.getByID"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", ctx.Get("request_id").(string)),
	)

	location, err := h.service.GetLocationById(int(id))
	if err != nil {
		log.Error("Failed to get location:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get location",
		})
	}

	log.Info("Location fetched successfully")

	return ctx.JSON(http.StatusOK, location)
}

func (h *LocationHandler) PostLocation(ctx echo.Context) error {
	const op = "rest.Location.create"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", ctx.Get("request_id").(string)),
	)

	var location locationapi.Location
	if err := decodeAndValidateEcho(ctx, &location, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.CreateLocation(location)
	if err != nil {
		log.Error("Failed to create location:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create location",
		})
	}

	log.Info("Location created successfully")

	return ctx.JSON(http.StatusCreated, resp)
}

func (h *LocationHandler) PutLocationId(ctx echo.Context, id locationapi.Id) error {
	const op = "rest.Location.update"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", ctx.Get("request_id").(string)),
	)

	var location locationapi.LocationUpdate
	if err := decodeAndValidateEcho(ctx, &location, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	resp, err := h.service.UpdateLocation(int(id), location)
	if err != nil {
		log.Error("Failed to update location:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update location",
		})
	}

	log.Info("Location updated successfully")

	return ctx.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) DeleteLocationId(ctx echo.Context, id locationapi.Id) error {
	const op = "rest.Location.delete"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", ctx.Get("request_id").(string)),
	)

	err := h.service.DeleteLocation(int(id))
	if err != nil {
		log.Error("Failed to delete location:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete location",
		})
	}

	log.Info("Location deleted successfully")

	return ctx.NoContent(http.StatusOK)
}

func NewLocation(log *slog.Logger, service *service.LocationService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	handler := NewLocationHandler(log, service, validate)

	r.Route("/location", func(r chi.Router) {
		r.Get("/", HandlerAdapter(handler.GetLocation))
		r.Post("/", HandlerAdapter(handler.PostLocation))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("id"))
				return handler.GetLocationId(ctx, locationapi.Id(id))
			}))

			r.Put("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("id"))
				return handler.PutLocationId(ctx, locationapi.Id(id))
			}))

			r.Delete("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("id"))
				return handler.DeleteLocationId(ctx, locationapi.Id(id))
			}))
		})
	})

	return r
}
