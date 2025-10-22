package models

type User struct {
	ID      string `json:"id" validate:"omitempty,uuid4"`
	Name    string `json:"name" validate:"required,min=1,max=60"`
	Age     int    `json:"age" validate:"required,min=1,max=100"`
	Country string `json:"country" validate:"required,alpha,len=2"`
}
