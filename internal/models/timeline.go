package models

import "time"

type Timeline struct {
	tableName struct{} `pg:"timeline"`

	ID          int       `pg:"id,pk"`
	Title       string    `pg:"title,type:varchar(255),notnull"`
	Description string    `pg:"description"`
	Deadline    time.Time `pg:"deadline"`
	IsBlocking  bool      `pg:"is_blocking,notnull"`
	IsScoring   bool      `pg:"is_scoring"`
	Status      string    `pg:"status"`

	TrackID int    `pg:"track_id"`
	Track   *Track `pg:"rel:has-one"`

	TimelineStatusID int             `pg:"timeline_status_id"`
	TimeLineStatus   *TimelineStatus `pg:"rel:has-one"`

	TimesActionStatuses []TeamActionStatus `pg:"many2many:team_action_status"`
}
