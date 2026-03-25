package dto

type UsersFilterRequest struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	PhoneNumber string   `json:"phoneNumber"`
	Roles       []string `json:"roles"`
}
