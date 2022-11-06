package strava

import (
	"errors"

	"golang.org/x/oauth2"

	config "github.com/tommzn/go-config"
	secrets "github.com/tommzn/go-secrets"
	strava "github.com/tommzn/go-strava"
)

// NewTokenSource creates a new token source which handles initial token request and token refresh.
// STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET is mandatory and have to be provided by used secrets manager.
// If STRAVA_REFRESH_TOKEN is avaiable via secrets manager it will directly be used to obtain new access token.
// Aithout a refresh token a authorization code (STRAVA_AUTH_CODE) is required to obtain access and refresh tokens.
func newTokenSource(conf config.Config, secretsmanager secrets.SecretsManager) (oauth2.TokenSource, error) {

	tokenUrl := conf.Get("strava.tokenurl", nil)
	if tokenUrl == nil {
		return nil, errors.New("No token url found in config.")
	}

	clientId, err := secretsmanager.Obtain("STRAVA_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := secretsmanager.Obtain("STRAVA_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	oauthConfig := oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     *clientId,
		ClientSecret: *clientSecret,
		Scopes:       []string{"read", "activity:read"},
		Endpoint: oauth2.Endpoint{
			TokenURL:  *tokenUrl,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	if refreshTokenStr, err := secretsmanager.Obtain("STRAVA_REFRESH_TOKEN"); err == nil {
		return strava.TokenSourceFromRefreshToken(oauthConfig, *refreshTokenStr)
	}

	authCode, err := secretsmanager.Obtain("STRAVA_AUTH_CODE")
	if err != nil {
		return nil, err
	}
	return strava.TokenSourceFromAuthorizationCode(oauthConfig, *authCode)
}
