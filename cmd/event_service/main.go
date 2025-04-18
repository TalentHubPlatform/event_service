package main

import (
	"event_service/internal/config"
	"event_service/internal/models"
	"event_service/internal/repositories"
	"event_service/internal/routes/rest"
	"event_service/internal/service"
	"event_service/pkg/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	logger.Info("Starting event-service", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	db := pg.Connect(&pg.Options{
		Addr:     cfg.SQLDatabase.Addr,
		User:     cfg.SQLDatabase.User,
		Password: cfg.SQLDatabase.Password,
		Database: cfg.SQLDatabase.Database,
	})

	InitM2M()
	InitPrometheus()

	router := chi.NewRouter()

	router.Mount("/event", createEventHandler(db, logger))
	router.Mount("/dates", createDateHandler(db, logger))
	router.Mount("/status", createStatusHandler(db, logger))
	router.Mount("/location", createLocationHandler(db, logger))
	router.Mount("/track", createTrackHandler(db, logger))
	router.Mount("/timeline", createTimelineHandler(db, logger))
	router.Mount("/team-action-status", createTeamActionStatusHandler(db, logger))
	router.Mount("/track-winner", createTrackWinnerHandler(db, logger))

	router.Handle("/metrics", promhttp.Handler())

	logger.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", err.Error())
	}

	logger.Error("server stopped")
}

func InitM2M() {
	orm.RegisterTable((*models.EventLocation)(nil))
	orm.RegisterTable((*models.EventPrize)(nil))
	orm.RegisterTable((*models.TeamActionStatus)(nil))
	orm.RegisterTable((*models.LocationTrack)(nil))
	orm.RegisterTable((*models.TrackJudge)(nil))
	orm.RegisterTable((*models.TrackWinner)(nil))
}

func InitPrometheus() {
	prometheus.MustRegister(utils.CronTaskSuccess)
	prometheus.MustRegister(utils.CronTaskFailure)
	prometheus.MustRegister(utils.EventStartSuccess)
	prometheus.MustRegister(utils.EventStartFailure)
	prometheus.MustRegister(utils.EventEndSuccess)
	prometheus.MustRegister(utils.EventEndFailure)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func createDateHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	dateRepository := repositories.NewDateRepository(db)
	dateService := service.NewDateService(dateRepository, db)

	return rest.NewDate(logger, dateService)
}

func createEventHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	eventRepository := repositories.NewEventRepository(db)
	eventLocationRepository := repositories.NewEventLocationRepository(db)
	trackRepository := repositories.NewTrackRepository(db)

	eventService := service.NewEventsService(eventRepository, trackRepository, eventLocationRepository, db)
	go utils.ScheduleEvents(logger, eventService)

	return rest.NewEvent(logger, eventService)
}

func createStatusHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	statusRepository := repositories.NewStatusRepository(db)
	statusService := service.NewStatusService(statusRepository, db)

	return rest.NewStatus(logger, statusService)
}

func createLocationHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	locationRepository := repositories.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepository, db)

	return rest.NewLocation(logger, locationService)
}

func createTrackHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	trackRepository := repositories.NewTrackRepository(db)

	trackService := service.NewTrackService(trackRepository, db)
	go utils.ScheduleTracks(logger, trackService)

	return rest.NewTrack(logger, trackService)
}

func createTimelineHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	timelineRepository := repositories.NewTimelineRepository(db)
	timelineStatusRepository := repositories.NewTimelineStatusRepository(db)

	timelineService := service.NewTimelineService(timelineRepository, timelineStatusRepository, db)
	return rest.NewTimeline(logger, timelineService)
}

func createTeamActionStatusHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	teamActionStatusRepository := repositories.NewTeamActionStatusRepository(db)
	teamActionStatusService := service.NewTeamActionStatusService(teamActionStatusRepository, db)

	return rest.NewTeamActionStatus(logger, teamActionStatusService)
}

func createTrackWinnerHandler(db *pg.DB, logger *slog.Logger) *chi.Mux {
	trackWinnerRepository := repositories.NewTrackWinnerRepository(db)
	trackRepository := repositories.NewTrackRepository(db)
	timelineRepository := repositories.NewTimelineRepository(db)

	trackWinnerService := service.NewTrackWinnerService(trackWinnerRepository, trackRepository, timelineRepository, db)
	return rest.NewTrackWinner(logger, trackWinnerService)
}
