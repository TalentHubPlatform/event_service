package models

import "time"

type Event struct {
	tableName    struct{}  `pg:"event"`
	ID           int       `pg:"id,pk"`
	Title        string    `pg:"title,type:varchar(255),unique,notnull"`
	Description  string    `pg:"description"`
	RedirectLink string    `pg:"redirect_link,notnull"`
	CreatedAt    time.Time `pg:"created_at,default:now()"`
	UpdatedAt    time.Time `pg:"updated_at"`
	Status       string    `pg:"status"`

	DateID int   `pg:"date_id"`
	Date   *Date `pg:"rel:has-one"`

	Locations   []Location   `pg:"many2many:event_location"`
	Tracks      []Track      `pg:"rel:has-many"`
	EventPrizes []EventPrize `pg:"rel:has-many"`
}
