package model

type RegisterModel struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password,omitempty" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password,omitempty" binding:"required,min=8"`
}
