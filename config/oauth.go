package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig *oauth2.Config
)

func InitGoogleOauthConfig() {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientId == "" || clientSecret == "" {
		logrus.Fatal("Please set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET env variable")
	}

	GoogleOauthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	logrus.Infof("Google Oauth initialized with client ID: %s", clientId)
}
