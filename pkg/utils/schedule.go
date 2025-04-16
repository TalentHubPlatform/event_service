package utils

import (
	"event_service/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"log/slog"
	"strconv"
)

var (
	CronTaskSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cron_task_success_total",
		Help: "Total number of successful cron tasks",
	})
	CronTaskFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cron_task_failure_total",
		Help: "Total number of failed cron tasks",
	})

	EventStartSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "event_start_success_total",
		Help: "Total number of successful event starts",
	})
	EventStartFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "event_start_failure_total",
		Help: "Total number of failed event starts",
	})
	EventEndSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "event_end_success_total",
		Help: "Total number of successful event ends",
	})
	EventEndFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "event_end_failure_total",
		Help: "Total number of failed event ends",
	})

	TrackStartSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "track_start_success_total",
		Help: "Total number of successful track starts",
	})
	TrackStartFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "track_start_failure_total",
		Help: "Total number of failed track starts",
	})
	TrackEndSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "track_end_success_total",
		Help: "Total number of successful track ends",
	})
	TrackEndFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "track_end_failure_total",
		Help: "Total number of failed track ends",
	})
)

type EventService interface {
	GetAllEventsToStart() ([]*models.Event, error)
	GetAllEventsToEnd() ([]*models.Event, error)

	StartEvent(int) (*models.Event, error)
	EndEvent(int) (*models.Event, error)
}

type TrackService interface {
	GetAllTracksToStart() ([]*models.Track, error)
	GetAllTracksToEnd() ([]*models.Track, error)

	StartTrack(int) (*models.Track, error)
	EndTrack(int) (*models.Track, error)
}

func ScheduleEvents(log *slog.Logger, service EventService) {
	c := cron.New()

	_, err := c.AddFunc("@every 1m", func() {
		events, err := service.GetAllEventsToStart()
		if err != nil {
			log.Error("Failed to get events to start")
			return
		}

		for _, event := range events {
			_, err = service.StartEvent(event.ID)

			if err != nil {
				log.Error("Failed to start event", slog.String("id", strconv.Itoa(event.ID)))
				EventStartFailure.Inc()
			} else {
				log.Info("Successfully started event", slog.String("id", strconv.Itoa(event.ID)))
				EventStartSuccess.Inc()
			}
		}
	})

	if err != nil {
		log.Error("Error scheduling start events", slog.String("error", err.Error()))
		CronTaskFailure.Inc()
	} else {
		CronTaskSuccess.Inc()
	}

	_, err = c.AddFunc("@every 1m", func() {
		events, err := service.GetAllEventsToEnd()
		if err != nil {
			log.Error("Failed to get events to end")
			return
		}

		for _, event := range events {
			_, err = service.EndEvent(event.ID)

			if err != nil {
				log.Error("Failed to env event", slog.String("id", strconv.Itoa(event.ID)))
				EventEndFailure.Inc()
			} else {
				log.Info("Successfully ended event", slog.String("id", strconv.Itoa(event.ID)))
				EventEndSuccess.Inc()
			}
		}
	})

	if err != nil {
		log.Error("Error scheduling end events:", slog.String("error", err.Error()))
		CronTaskFailure.Inc()
	} else {
		CronTaskSuccess.Inc()
	}

	c.Start()
}

func ScheduleTracks(log *slog.Logger, service TrackService) {
	c := cron.New()

	_, err := c.AddFunc("@every 1m", func() {
		tracks, err := service.GetAllTracksToStart()
		if err != nil {
			log.Error("Failed to get tracks to start")
			return
		}

		for _, track := range tracks {
			_, err = service.StartTrack(track.ID)

			if err != nil {
				log.Error("Failed to start track", slog.String("id", strconv.Itoa(track.ID)))
				TrackStartFailure.Inc()
			} else {
				log.Info("Successfully started track", slog.String("id", strconv.Itoa(track.ID)))
				TrackStartSuccess.Inc()
			}
		}
	})

	if err != nil {
		log.Error("Error scheduling start tracks", slog.String("error", err.Error()))
		CronTaskFailure.Inc()
	} else {
		CronTaskSuccess.Inc()
	}

	_, err = c.AddFunc("@every 1m", func() {
		tracks, err := service.GetAllTracksToEnd()
		if err != nil {
			log.Error("Failed to get tracks to end")
			return
		}

		for _, track := range tracks {
			_, err = service.EndTrack(track.ID)

			if err != nil {
				log.Error("Failed to env track", slog.String("id", strconv.Itoa(track.ID)))
				TrackEndFailure.Inc()
			} else {
				log.Info("Successfully ended track", slog.String("id", strconv.Itoa(track.ID)))
				TrackEndSuccess.Inc()
			}
		}
	})

	if err != nil {
		log.Error("Error scheduling end tracks:", slog.String("error", err.Error()))
		CronTaskFailure.Inc()
	} else {
		CronTaskSuccess.Inc()
	}

	c.Start()
}
