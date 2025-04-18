package schemas

import "time"

type Timeline struct {
	Title       string    `json:"title" validate:"required" example:"Timeline"`
	Description string    `json:"description" validate:"required" example:"Description"`
	Deadline    time.Time `json:"deadline" validate:"required" example:"2020-01-01 00:00:00"`
	IsBlocking  bool      `json:"is_blocking" validate:"omitempty" example:"true"`
	Status      string    `json:"status" validate:"required" example:"ready"`
	TrackID     int       `json:"track_id" validate:"required" example:"1"`
}

type TimelineUpdate struct {
	Title            string    `json:"title" example:"Timeline"`
	Description      string    `json:"description" example:"Description"`
	Deadline         time.Time `json:"deadline" example:"2020-01-01 00:00:00"`
	IsBlocking       bool      `json:"is_blocking" validate:"required" example:"true"`
	Status           string    `json:"status" example:"ready"`
	TrackID          int       `json:"track_id" example:"1"`
	TimelineStatusID int       `json:"timeline_status_id" example:"1"`
}
