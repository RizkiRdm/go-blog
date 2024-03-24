package blog

type RequestCreateBlog struct {
	UsernameId string   `json:"user_id"`
	Thumbnail  string   `json:"thumbnail"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	CategoryId string   `json:"category_id"`
	TagsId     []string `json:"tags_id"`
}

type RequestUpdateBlog struct {
	Thumbnail  string   `json:"thumbnail"`
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	CategoryId string   `json:"category_id"`
	TagsId     []string `json:"tags_id"`
}
