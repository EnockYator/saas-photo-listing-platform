package oidc

import (
	"context"
	"fmt"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/domain"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewProvider(
	ctx context.Context,
	cfg Config,
) (*Provider, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(
		&oidc.Config{
			ClientID: cfg.ClientID,
		},
	)

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			oidc.ScopeProfile,
			oidc.ScopeEmail,
		},
	}

	return &Provider{
		config:   oauthConfig,
		verifier: verifier,
	}, nil
}

func (p *Provider) AuthorizationURL(
	state string,
	nonce string,
	codeChallenge string,
) (string, error) {
	url := p.config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(
			codeChallenge,
		),
		oauth2.SetAuthURLParam("nonce", nonce),
	)

	return url, nil
}

func (p *Provider) ExchangeCode(
	ctx context.Context,
	code string,
	codeVerifier string,
) (*domain.IdentityClaims, error) {
	token, err := p.config.Exchange(
		ctx,
		code,
		oauth2.VerifierOption(codeVerifier),
	)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("%v", domain.ErrMissingIDToken)
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims struct {
		Subject       string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &domain.IdentityClaims{
		Issuer:        idToken.Issuer,
		Subject:       claims.Subject,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Name:          claims.Name,
	}, nil
}
