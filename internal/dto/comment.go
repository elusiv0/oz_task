package dto

import "time"

type Comment struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	ArticleID int       `json:"articleId"`
	ParentID  *int      `json:"parentId"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewComment struct {
	Text      string `json:"text"`
	ArticleID int    `json:"articleId"`
	ParentID  *int   `json:"parentId,omitempty"`
}

type GetCommentsRequest struct {
	PostId   *int `json:"post_id"`
	ParentId *int `json:"parent_id"`
	After    *int `json:"after"`
	First    int  `json:"first"`
}
