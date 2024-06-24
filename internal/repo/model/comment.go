package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id        int
	Text      string
	ArticleID int
	ParentId  sql.NullInt32
	CreatedAt time.Time
	Rown      *int
}
