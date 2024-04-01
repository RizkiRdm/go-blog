package category

import "time"

type categories struct {
	Id        uint      `json:"id"`
	Name      string    `json:"category_name"`
	CreatedAt time.Time `json:"created_at"`
}

type categoryResponse struct {
	Name      string    `json:"category_name"`
	CreatedAt time.Time `json:"created_at"`
}
