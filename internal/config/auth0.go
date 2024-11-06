package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func NewAuth0(config *viper.Viper) (*Authenticator, error) {
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.GetString("auth0.domain")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     config.GetString("auth0.client_id"),
		ClientSecret: config.GetString("auth0.client_secret"),
		RedirectURL:  config.GetString("auth0.redirect_url"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
