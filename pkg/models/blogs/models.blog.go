package blogs

import "time"

type Blogs struct {
	Id         uint      `json:"id"`
	UserId     uint      `json:"user_id"`
	CategoryId uint      `json:"category_id"`
	TagId      uint      `json:"tag_id"`
	Thumbnail  string    `json:"thumbnail"`
	Title      string    `json:"title_blog"`
	Body       string    `json:"body_blog"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
