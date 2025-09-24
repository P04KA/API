package models

type User struct {
	ID      string `json:"id" validate:"omitempty,min=1,max=100"`
	Name    string `json:"name" validate:"omitempty,min=1,max=100"`
	Age     int    `json:"age" validate:"omitempty,min=1,max=150"`
	Country string `json:"country" validate:"omitempty,min=1,max=2"`
}
