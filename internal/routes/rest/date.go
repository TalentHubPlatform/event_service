package rest

import (
	api "event_service/gen/date"
	"event_service/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
)

type DateService interface {
	GetAllDates() ([]*api.DateResponse, error)
	GetDateByID(id int) (*api.DateResponse, error)
	CreateDate(date api.Date) (*api.DateResponse, error)
	UpdateDate(id int, date api.DateUpdate) (*api.DateResponse, error)
	DeleteDate(id int) error
}

type DateHandler struct {
	log       *slog.Logger
	service   DateService
	validator *validator.Validate
}

func NewDateHandler(log *slog.Logger, service DateService, validator *validator.Validate) *DateHandler {
	return &DateHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func decodeAndValidateEcho(ctx echo.Context, dst interface{}, validate *validator.Validate) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	if err := ctx.Bind(dst); err != nil {
		return fmt.Errorf("failed to bind request body: %w", err)
	}

	return nil
}

func (h *DateHandler) GetDates(ctx echo.Context) error {
	const op = "rest.Date.getAll"

	log := slog.With(
		slog.String("op", op),
	)

	dates, err := h.service.GetAllDates()
	if err != nil {
		log.Error("Failed to get dates:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get dates",
		})
	}

	log.Info("Dates fetched successfully")

	return ctx.JSON(http.StatusOK, dates)
}

func (h *DateHandler) PostDates(ctx echo.Context) error {
	const op = "rest.Date.create"

	log := h.log.With(
		slog.String("op", op),
	)

	var date api.Date
	if err := decodeAndValidateEcho(ctx, &date, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.CreateDate(date)
	if err != nil {
		log.Error("Failed to create date:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create date"})
	}

	log.Info("Date created successfully")

	return ctx.JSON(http.StatusCreated, resp)
}

func (h *DateHandler) PutDatesId(ctx echo.Context, id api.Id) error {
	const op = "rest.Date.update"

	log := h.log.With(
		slog.String("op", op),
	)

	var date api.DateUpdate
	if err := decodeAndValidateEcho(ctx, &date, h.validator); err != nil {
		log.Error("Failed to decode and validate request:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.service.UpdateDate(int(id), date)
	if err != nil {
		log.Error("Failed to update date:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update date"})
	}

	log.Info("Date updated successfully")

	return ctx.JSON(http.StatusOK, resp)
}

func (h *DateHandler) GetDatesId(ctx echo.Context, id api.Id) error {
	const op = "rest.Date.getByID"

	log := h.log.With(
		slog.String("op", op),
	)

	date, err := h.service.GetDateByID(int(id))
	if err != nil {
		log.Error("Failed to get date:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get date"})
	}

	log.Info("Date fetched successfully")

	return ctx.JSON(http.StatusOK, date)
}

func (h *DateHandler) DeleteDatesId(ctx echo.Context, id api.Id) error {
	const op = "rest.Date.delete"

	log := h.log.With(
		slog.String("op", op),
	)

	err := h.service.DeleteDate(id)
	if err != nil {
		log.Error("Failed to delete date:", slog.String("error", err.Error()))

		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete date"})
	}

	log.Info("Date deleted successfully")

	return ctx.NoContent(http.StatusOK)
}

type echoResponseWriter struct {
	http.ResponseWriter
}

func (w *echoResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func HandlerAdapter(echoHandler func(echo.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		echoContext := echo.New().NewContext(r, &echoResponseWriter{w})

		params := chi.RouteContext(r.Context()).URLParams
		for idx, key := range params.Keys {
			echoContext.SetParamNames(key)
			echoContext.SetParamValues(params.Values[idx])
		}

		_ = echoHandler(echoContext)
	}
}

func NewDate(log *slog.Logger, service *service.DateService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.Validate{}

	handler := NewDateHandler(log, service, &validate)

	r.Route("/", func(r chi.Router) {
		r.Get("/", HandlerAdapter(handler.GetDates))
		r.Post("/", HandlerAdapter(handler.PostDates))

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				return handler.GetDatesId(ctx, api.Id(id))
			}))

			r.Put("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				return handler.PutDatesId(ctx, api.Id(id))
			}))

			r.Delete("/", HandlerAdapter(func(ctx echo.Context) error {
				id, _ := strconv.Atoi(ctx.Param("Id"))
				return handler.DeleteDatesId(ctx, api.Id(id))
			}))
		})
	})

	return r
}
