package model

import "time"

type Post struct {
	Id        int
	Title     string
	Text      string
	Closed    bool
	CreatedAt time.Time
}
