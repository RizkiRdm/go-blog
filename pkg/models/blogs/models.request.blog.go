package blogs

type RequestCreateBlog struct {
	UserId     uint   `json:"user_id"`
	CategoryId uint   `json:"category_id"`
	TagId      uint   `json:"tag_id"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title_blog"`
	Body       string `json:"body_blog"`
}

type RequestUpdateBlog struct {
	CategoryId uint   `json:"category_id"`
	TagId      uint   `json:"tag_id"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title_blog"`
	Body       string `json:"body_blog"`
}
