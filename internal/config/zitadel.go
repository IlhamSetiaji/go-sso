package config

import (
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type ZitadelAuthenticator struct {
	oauth2.Config
}

type ZitadelAppConfig struct {
	RedirectURI string
}

var ZitadelAppConfigs = map[string]ZitadelAppConfig{
	"google": {RedirectURI: "https://www.google.com"},
	"github": {RedirectURI: "https://www.github.com"},
}

func NewZitadelAuthenticator(config *viper.Viper) (*ZitadelAuthenticator, error) {
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	conf := oauth2.Config{
		ClientID:     config.GetString("zitadel.client_id"),
		ClientSecret: config.GetString("zitadel.client_secret"),
		RedirectURL:  config.GetString("zitadel.redirect_url"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GetString("zitadel.auth_url"),  // Zitadel Auth URL
			TokenURL: config.GetString("zitadel.token_url"), // Zitadel Token URL
		},
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &ZitadelAuthenticator{
		Config: conf,
	}, nil
}
