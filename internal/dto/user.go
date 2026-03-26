package dto

type UserRequest struct {
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password"`
	FirstName   string   `json:"firstName" validate:"required"`
	LastName    string   `json:"lastName" validate:"required"`
	PhoneNumber string   `json:"phoneNumber" validate:"phoneNumber"`
	Avatar      string   `json:"avatar"`
	Nip         string   `json:"nip"`
	Roles       []string `json:"roles"`
}
