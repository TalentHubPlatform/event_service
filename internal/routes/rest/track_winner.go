package rest

import (
	"encoding/json"
	"event_service/internal/models"
	"event_service/internal/repositories"
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

type TrackWinnerService interface {
	GetWinnersOfTrack(int) ([]*models.TrackWinner, error)
	GetWinnerById(int, int) (*models.TrackWinner, error)
	CreateWinnerOfTrack(*schemas.TrackWinner) (*models.TrackWinner, error)

	CalculateRating(int, int, int) ([]*repositories.AggregateResult, error)
}

func NewTrackWinner(log *slog.Logger, service *service.TrackWinnerService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Logger)

	validate := validator.New()

	r.Route("/", func(r chi.Router) {
		r.Get("/", getWinnersOfTrackHandler(log, service))
		r.Post("/", createWinnersOfTrackHandler(log, service, validate))

		r.Route("/{trackId}", func(r chi.Router) {
			r.Get("/", getWinnerByIdHandler(log, service))
			r.Post("/", calculateRatingHandler(log, service))
		})
	})

	return r
}

func getWinnersOfTrackHandler(log *slog.Logger, service TrackWinnerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TrackWinners.getAllOfTrack"

		log := log.With(
			slog.With("op", op),
			slog.With("request_id", middleware.GetReqID(r.Context())),
		)

		headersList := map[string]string{
			"TeamId": "int",
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Headers validation failed with error:", err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tracks, err := service.GetWinnersOfTrack(convertedHeaders["TeamId"].(int))
		if err != nil {
			log.Error("error getting winners:", err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tracks); err != nil {
			log.Error("error encoding tracks:", err)

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("Winners successfully fetched")
	}
}

func createWinnersOfTrackHandler(log *slog.Logger, service TrackWinnerService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TrackWinners.create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var track schemas.TrackWinner
		if err := DecodeAndValidate(r, &track, validate); err != nil {
			log.Error("Failed to decode request:", err.Error())

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateWinnerOfTrack(&track)
		if err != nil {
			log.Error("Failed to create TrackWinner:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TrackWinner created successfully")
	}
}

func getWinnerByIdHandler(log *slog.Logger, service TrackWinnerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TrackWinner.getByID"

		log := log.With(
			slog.String("op", op),
			slog.String("request_it", middleware.GetReqID(r.Context())),
		)

		trackId, err := strconv.Atoi(chi.URLParam(r, "trackId"))

		headersList := map[string]string{
			"TeamId": "int",
		}

		convertedHeaders, err := utils.ValidateHeaders(headersList, log, r)
		if err != nil {
			log.Error("Headers validation failed with error: ", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		track, err := service.GetWinnerById(trackId, convertedHeaders["TeamId"].(int))
		if err != nil {
			log.Error("Failed to get TrackWinner:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(track); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TrackWinner fetched successfully")
	}
}

func calculateRatingHandler(log *slog.Logger, service TrackWinnerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.TrackWinner.getByID"

		log := log.With(
			slog.String("op", op),
			slog.String("request_it", middleware.GetReqID(r.Context())),
		)

		trackId, _ := strconv.Atoi(chi.URLParam(r, "trackId"))
		queryParams := r.URL.Query()

		limit := 100
		offset := 0

		if queryParams.Get("limit") != "" {
			converted, err := strconv.Atoi(queryParams.Get("status_id"))
			if err != nil {
				log.Error("Failed to convert limit query param:", err.Error())

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			limit = converted
		}

		if queryParams.Get("offset") != "" {
			converted, err := strconv.Atoi(queryParams.Get("offset"))
			if err != nil {
				log.Error("Failed to convert offset query param:", err.Error())

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			offset = converted
		}

		results, err := service.CalculateRating(trackId, limit, offset)
		if err != nil {
			log.Error("Failed to calculate results:", err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Error("Failed to encode response:", err.Error())

			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

		log.Info("TrackWinners results calculated successfully")
	}
}
