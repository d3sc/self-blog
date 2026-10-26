package models

import "time"

type Admin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Excerpt     string     `json:"excerpt"`
	Content     string     `json:"content,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	PublishedAt *time.Time `json:"publishedAt"`
}
type Page struct {
	Items      []Post     `json:"items"`
	Pagination Pagination `json:"pagination"`
}
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
type PostInput struct {
	Title   string  `json:"title"`
	Slug    *string `json:"slug"`
	Excerpt string  `json:"excerpt"`
	Content string  `json:"content"`
	Status  string  `json:"status"`
}
