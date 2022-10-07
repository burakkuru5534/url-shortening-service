package model

type RegisterRequest struct {
	ID        int64  `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	IsPremium bool   `json:"IsPremium"`
}

func RegisterRequestData() *RegisterRequest {
	m := new(RegisterRequest)

	return m
}
