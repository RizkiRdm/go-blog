package category

type RequestCreateCategory struct {
	Name string `json:"category_name"`
}

type RequestUpdateCategory struct {
	Name string `json:"category_name"`
}
