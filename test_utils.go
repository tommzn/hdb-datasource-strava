package strava

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"golang.org/x/oauth2"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
)

func loadConfigForTest(configFile string) config.Config {
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

func secretsManagerForTest() secrets.SecretsManager {
	secretsMap := make(map[string]string)
	secretsMap["STRAVA_ATHLETE_ID"] = "123"
	secretsMap["STRAVA_CLIENT_ID"] = "123456789"
	secretsMap["STRAVA_CLIENT_SECRET"] = "987654321"
	secretsMap["STRAVA_REFRESH_TOKEN"] = "xxx"
	return secrets.NewStaticSecretsManager(secretsMap)
}

// loggerForTest creates a new stdout logger for testing.
func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func serverForTest() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizedAthleteRegex, _ := regexp.Compile("/athlete$")
		activitiesRegex, _ := regexp.Compile("/athlete/activities$")
		statsRegex, _ := regexp.Compile("/athletes/[0-9]{1,}/stats$")
		switch true {
		case authorizedAthleteRegex.MatchString(r.URL.Path):
			serveFile(w, r, "fixtures/authorizedAthlete.json")
		case activitiesRegex.MatchString(r.URL.Path):
			serveFile(w, r, "fixtures/athleteactivities.json")
		case statsRegex.MatchString(r.URL.Path):
			serveFile(w, r, "fixtures/athletestats.json")
		default:
			serveError(w, r)
		}
	}))
}

func serverWithErrorResponseForTest() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveError(w, r)
	}))
}

func serveFile(w http.ResponseWriter, r *http.Request, filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func serveError(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("fixtures/fault.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

// NewTokenSourceMock returns a local token source for testing.
// If shouldReturnError is true Token() method will always return with an error.
func newTokenSourceMock(shouldReturnError bool) *tokenSourceMock {
	return &tokenSourceMock{shouldReturnError: shouldReturnError}
}

// TokenSourceMock is a local source of tokens, to be used for testing.
type tokenSourceMock struct {
	shouldReturnError bool
}

// Token returns a dummy token for testing.
// If shouldReturnError is true an error will be returned.
func (mock *tokenSourceMock) Token() (*oauth2.Token, error) {
	if mock.shouldReturnError {
		return nil, errors.New("An Error has occurred.")
	}
	return &oauth2.Token{
		AccessToken:  "<ACCESS_TOKEN>",
		TokenType:    "Bearer",
		RefreshToken: "<REFRESH_TOKEN>",
		Expiry:       time.Now().Add(30 * time.Minute),
	}, nil
}
