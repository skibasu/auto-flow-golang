package models

type User struct {
	Id          string   `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	PhoneNumber string   `json:"phoneNumber"`
	Avatar      string   `json:"avatar"`
	Roles       []string `json:"roles"`
}
