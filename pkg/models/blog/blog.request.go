package blog

type RequestCreateBlog struct {
	UsernameId int    `json:"user_id"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryId int    `json:"category_id"`
}

type RequestUpdateBlog struct {
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryId string `json:"category_id"`
}
