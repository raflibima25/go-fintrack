package response

import "time"

type LoginResponse struct {
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
	Expiration  time.Time `json:"expiration"`
	IsAdmin     bool      `json:"is_admin"`
}
