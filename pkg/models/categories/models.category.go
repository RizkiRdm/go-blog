package models

import "time"

type Categories struct {
	Id         uint      `json:"id"`
	Title      string    `json:"category_name"`
	Created_at time.Time `json:"Created_at"`
}
