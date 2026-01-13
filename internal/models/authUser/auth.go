package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` //will never be encoded
	CreatedAt time.Time `json:"created_at"`
}
type AuthReqBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"` //will never be encoded
}

type UserResponse struct {
	ID        int       `json:"id" validate:"required,gt=0"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type ChangePasswordUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
