package schemas

import "time"

type TeamActionStatus struct {
	TrackTeamID    int       `json:"track_team_id" validate:"required" example:"1"`
	TimelineID     int       `json:"timeline_id" validate:"required" example:"1"`
	ResultValue    int       `json:"result_value" validate:"required" example:"600"`
	ResolutionLink string    `json:"resolution_link" validate:"required" example:"https://www.youtube.com"`
	CompletedAt    time.Time `json:"completed_at" validate:"required" example:"2023-01-02T00:00:00Z"`
	Notes          string    `json:"notes" validate:"required" example:"Notes"`
}

type TeamActionStatusUpdate struct {
	ResultValue    string    `json:"result_value" example:"600"`
	ResolutionLink string    `json:"resolution_link" example:"https://www.youtube.com"`
	CompletedAt    time.Time `json:"completed_at" example:"2023-01-02T00:00:00Z"`
	Notes          string    `json:"notes" example:"Notes"`
}
