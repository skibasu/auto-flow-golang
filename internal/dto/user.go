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

type UpdateUserRequest struct {
	Password    *string  `json:"password" validate:"omitempty,min=8"`
	Email       *string  `json:"email" validate:"omitempty,email"`
	PhoneNumber *string  `json:"phoneNumber" validate:"omitempty,phoneNumber"`
	FirstName   *string  `json:"firstName" validate:"omitempty,required,min=3,max=36"`
	LastName    *string  `json:"lastName" validate:"omitempty,required,min=3,max=36"`
	Nip         *string  `json:"nip" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
}
