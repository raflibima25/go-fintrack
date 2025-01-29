package request

type RegisterRequest struct {
	Name            string `json:"name" binding:"required,min=3,max=50"`
	Username        string `json:"username" binding:"required,alphanum,min=4,max=20"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password,omitempty" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password,omitempty" binding:"required,eqfield=Password"`
}
