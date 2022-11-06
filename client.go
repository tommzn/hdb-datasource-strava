package strava

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	secrets "github.com/tommzn/go-secrets"
	strava "github.com/tommzn/go-strava"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// New create a Strava client.
// Required config is:
//	strava.athleteid:		Id of an athelte at Strave this client should retrieve stats and activities for.
//	strava.activitycount:	Number of retrieved activities. Default is 5 activities.
//	strava.activitydays:	Time range of days activities should be retrieved for.
func New(conf config.Config, secretsmanager secrets.SecretsManager, logger log.Logger) (*StravaClient, error) {

	athleteIdStr := conf.Get("strava.athleteid", config.AsStringPtr(os.Getenv("STRAVA_ATHLETE_ID")))
	if athleteIdStr == nil {
		return nil, errors.New("No athlete Id found in config.")
	}
	athleteId, _ := strconv.ParseInt(*athleteIdStr, 10, 64)
	tokenSource, err := newTokenSource(conf, secretsmanager)
	if err != nil {
		return nil, err
	}
	activityCount := conf.GetAsInt("strava.activitycount", config.AsIntPtr(5))
	activityDays := conf.GetAsInt("strava.activitydays", config.AsIntPtr(30))
	stravaApiClient := strava.New(tokenSource)
	stravaApiClient.WithAthleteId(athleteId)
	return &StravaClient{
		stravaApiClient: stravaApiClient,
		activityCount:   *activityCount,
		activityDays:    *activityDays,
		logger:          logger,
	}, nil
}

// Fetch calls stats and activity APIs to collect data for current athlete.
func (client *StravaClient) Fetch() (proto.Message, error) {

	stats, err := client.stravaApiClient.AthleteStats()
	if err != nil {
		return nil, err
	}

	pagination := strava.NewPagination(1, client.activityCount)
	after := time.Now().Add(time.Duration(client.activityDays) * time.Hour * 24 * -1)
	timeFilter := strava.TimeFilter{After: &after}
	activities, err := client.stravaApiClient.AthleteActivities(&timeFilter, pagination)
	if err != nil {
		return nil, err
	}

	return asEvent(stats, activities), nil
}

// AsEvent converts given Strava data into an event.
func asEvent(stats *strava.ActivityStats, listOfActivities *[]strava.SummaryActivity) proto.Message {
	activityStats := &events.ActivityStats{
		BiggestRideDistance:  stats.BiggestRideDistance,
		RecentRideTotals:     asActivityTotal(stats.RecentRideTotals),
		RecentRunTotals:      asActivityTotal(stats.RecentRunTotals),
		YearToDateRideTotals: asActivityTotal(stats.YearToDateRideTotals),
		YearToDateRunTotals:  asActivityTotal(stats.YearToDateRunTotals),
		AllRideotals:         asActivityTotal(stats.AllRideotals),
		AllRunTotals:         asActivityTotal(stats.AllRunTotals),
	}
	activities := []*events.Activity{}
	for _, activity := range *listOfActivities {
		activities = append(activities, asActivity(activity))
	}
	return &events.AthleteStats{
		ActivityStats: activityStats,
		Activities:    activities,
		Timestamp:     asTimeStamp(time.Now()),
	}
}

// AsActivityTotal converts given activity totals to event data.
func asActivityTotal(activityTotal strava.ActivityTotal) *events.ActivityTotal {
	return &events.ActivityTotal{Count: activityTotal.Count, Distance: activityTotal.Distance, MovingTime: activityTotal.MovingTime}
}

// AsActivity converts given activitiy to an event.
func asActivity(activity strava.SummaryActivity) *events.Activity {
	return &events.Activity{
		Name:       activity.Name,
		Distance:   activity.Distance,
		MovingTime: activity.MovingTime,
		SportType:  activity.SportType,
		Timestamp:  asTimeStamp(activity.StartDateLocal),
	}
}

// asTimeStamp converts a unix epoch timestamp to a Protobuf timestamp.
func asTimeStamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
