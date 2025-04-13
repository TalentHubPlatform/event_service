package models

type TrackTeam struct {
	tableName struct{} `pg:"track_team"`
	ID        int      `pg:"id,pk"`
	TeamID    int      `pg:"team_id,notnull"`
	IsActive  bool     `pg:"is_active,notnull"`

	TrackID int    `pg:"track_id"`
	Track   *Track `pg:"rel:has-one"`

	ActionStatuses []TeamActionStatus `pg:"many2many:team_action_status"`
}
