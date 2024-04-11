package category

import "time"

type Categories struct {
	Id        uint      `json:"id"`
	Name      string    `json:"category_name"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryResponse struct {
	Name string `json:"category_name"`
}
