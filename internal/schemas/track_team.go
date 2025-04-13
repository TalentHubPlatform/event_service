package schemas

type TrackTeam struct {
	TeamID   int  `json:"team_id" validate:"required" example:"1"`
	IsActive bool `json:"is_active" validate:"required" example:"true"`
	TrackID  int  `json:"track_id" validate:"required" example:"1"`
}

type TrackTeamUpdate struct {
	IsActive string `json:"is_active" example:"false"`
}
