package blog

type RequestCreateBlog struct {
	UsernameId int    `json:"username"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title_blog"`
	Body       string `json:"body_blog"`
}

type RequestUpdateBlog struct {
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title_blog"`
	Body      string `json:"body_blog"`
}
