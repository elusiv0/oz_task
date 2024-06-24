package dto

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Closed    bool      `json:"closed"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewPost struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Closed bool   `json:"closed"`
}

type GetPostsRequest struct {
	First int  `json:"first"`
	After *int `json:"after"`
}
