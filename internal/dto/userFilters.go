package dto

type UsersFilterRequest struct {
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	Roles       []string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
