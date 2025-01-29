package request

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required" error:"Email or username is required"`
	Password        string `json:"password" binding:"required,min=8" error:"Password is required and must be at least 8 characters"`
}
