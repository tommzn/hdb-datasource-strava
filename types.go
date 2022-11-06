package strava

import (
	log "github.com/tommzn/go-log"
	strava "github.com/tommzn/go-strava"
)

// StravaClient handles request to Strava APIs.
type StravaClient struct {

	// Strava API Client
	stravaApiClient *strava.Client

	// Number of latest activities.
	activityCount int

	// Days for last activities.
	activityDays int

	// Logger
	logger log.Logger
}
