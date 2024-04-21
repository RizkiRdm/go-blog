package blog

import "time"

type Blogs struct {
	Id         uint      `json:"id"`
	UserId     int       `json:"user_id"`
	CategoryId int       `json:"category_id"`
	Thumbnail  string    `json:"thumbnail"`
	Title      string    `json:"title_blog"`
	Body       string    `json:"body_blog"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BlogResponse struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	Category  string    `json:"category"`
	Thumbnail string    `json:"thumbnail"`
	Title     string    `json:"title_blog"`
	Body      string    `json:"body_blog"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BlogSearchResponse struct {
	Title string `json:"title_blog"`
	Body  string `json:"body_blog"`
}
