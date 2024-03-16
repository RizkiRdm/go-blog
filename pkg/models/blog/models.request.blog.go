package blog

type RequestCreateBlog struct {
	UsernameId int      `json:"username"`
	Thumbnail  string   `json:"thumbnail"`
	Title      string   `json:"title_blog"`
	Body       string   `json:"body_blog"`
	Category   string   `json:"cateogory"`
	Tags       []string `json:"tags"`
}

type RequestUpdateBlog struct {
	Thumbnail string   `json:"thumbnail"`
	Title     string   `json:"title_blog"`
	Body      string   `json:"body_blog"`
	Category  string   `json:"cateogory"`
	Tags      []string `json:"tags"`
}
