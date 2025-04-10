package models

type TimelineStatus struct {
	tableName struct{} `pg:"timeline_status"`

	ID       int `pg:"id,pk"`
	CountNum int `pg:"count_num,notnull"`
}
