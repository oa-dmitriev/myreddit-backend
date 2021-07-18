package post

import (
	"myreddit/pkg/user"
)

type Post struct {
	Votes     []*user.User `json:"votes"`
	Score     int          `json:"score"`
	Id        string       `json:"id"`
	Text      string       `json:"text"`
	Title     string       `json:"title"`
	Category  string       `json:"category"`
	Author    *user.User   `json:"author,omitempty"`
	Comments  []*Comment   `json:"comments"`
	CreatedAt string       `json:"created"`
}

type Comment struct {
	Id        string     `json:"id"`
	Author    *user.User `json:"author"`
	Body      string     `json:"body"`
	CreatedAt string     `json:"created"`
}

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
