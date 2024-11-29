package config

import (
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthenticator struct {
	// *oidc.Provider
	oauth2.Config
}

type GoogleAppConfig struct {
	RedirectURI string
}

var GoogleAppConfigs = map[string]GoogleAppConfig{
	"google": {RedirectURI: "https://www.google.com"},
	"github": {RedirectURI: "https://www.github.com"},
}

func NewGoogleAuthenticator(config *viper.Viper) (*GoogleAuthenticator, error) {
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	conf := oauth2.Config{
		ClientID:     config.GetString("google.client_id"),
		ClientSecret: config.GetString("google.client_secret"),
		RedirectURL:  config.GetString("google.redirect_url"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	return &GoogleAuthenticator{
		Config: conf,
	}, nil
}
