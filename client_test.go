package strava

import (
	"testing"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	strava "github.com/tommzn/go-strava"
	events "github.com/tommzn/hdb-events-go"
)

type StravaClientTestSuite struct {
	suite.Suite
	conf   config.Config
	logger log.Logger
}

func TestStravaClientTestSuite(t *testing.T) {
	suite.Run(t, new(StravaClientTestSuite))
}

func (suite *StravaClientTestSuite) SetupTest() {
	suite.conf = loadConfigForTest("fixtures/testconfig.yml")
	suite.logger = loggerForTest()
}

func (suite *StravaClientTestSuite) TestFetch() {

	ts := serverForTest()
	defer ts.Close()

	client, err := New(suite.conf, secretsManagerForTest(), suite.logger)
	suite.Nil(err)
	suite.NotNil(client)
	stravaApiClient := strava.New(newTokenSourceMock(false))
	stravaApiClient.WithBaseUrl(ts.URL)
	client.stravaApiClient = stravaApiClient

	event, err := client.Fetch()
	suite.Nil(err)
	suite.NotNil(event)
	suite.Len(event.(*events.AthleteStats).Activities, 2)
	suite.True(event.(*events.AthleteStats).ActivityStats.RecentRideTotals.Distance > 0)
}

func (suite *StravaClientTestSuite) TestFetchWithError() {

	ts := serverWithErrorResponseForTest()
	defer ts.Close()

	client, err := New(suite.conf, secretsManagerForTest(), suite.logger)
	suite.Nil(err)
	suite.NotNil(client)
	stravaApiClient := strava.New(newTokenSourceMock(false))
	stravaApiClient.WithBaseUrl(ts.URL)
	client.stravaApiClient = stravaApiClient

	event, err := client.Fetch()
	suite.NotNil(err)
	suite.Nil(event)
}
