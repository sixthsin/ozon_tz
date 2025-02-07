package models

import "time"

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      string    `json:"authorId"`
	AllowComments bool      `json:"allowComments"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	ParentID  *string   `json:"parentId,omitempty"`
	AuthorID  string    `json:"authorId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
