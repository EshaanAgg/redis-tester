package internal

import (
	"fmt"
	"math"
	"strings"

	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/go-redis/redis"
)

type GeoPosition struct {
	key       string
	isInvalid bool
	latitude  float64
	longitude float64
}

type GeoPosTest struct {
	key       string
	locations []GeoPosition
}

func getGeoPositionsFromLocations(locations []GeoLocation) []GeoPosition {
	var geoPositions []GeoPosition
	for _, loc := range locations {
		geoPositions = append(geoPositions, GeoPosition{
			key:       loc.key,
			isInvalid: false,
			longitude: loc.longitude,
			latitude:  loc.latitude,
		})
	}
	return geoPositions
}

func NewGeoPosTest(key string, locations []GeoLocation) *GeoPosTest {
	return &GeoPosTest{
		key:       key,
		locations: getGeoPositionsFromLocations(locations),
	}
}

// Returns the command string for the GeoHashTest in the format "geoadd key longitude latitude"
func (t *GeoPosTest) command() string {
	var members []string
	for _, l := range t.locations {
		members = append(members, l.key)
	}
	return fmt.Sprintf("geopos %s %s", t.key, strings.Join(members, " "))
}

func (t *GeoPosTest) Run(client *redis.Client, logger *logger.Logger) error {
	logger.Infof("$ redis-cli %s", t.command())

	var memberNames []string
	for _, l := range t.locations {
		memberNames = append(memberNames, l.key)
	}
	resp, err := client.GeoPos(t.key, memberNames...).Result()

	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	// Check the response
	if len(resp) != len(t.locations) {
		return fmt.Errorf("Expected %d positions, got %d", len(t.locations), len(resp))
	}

	for i, l := range t.locations {
		if l.isInvalid {
			if resp[i] != nil {
				return fmt.Errorf("Expected position for %q to be nil, got %v", l.key, resp[i])
			}
			logger.Successf("Correctly received nil position for %q", l.key)
		} else {
			pos := resp[i]
			longError := math.Abs(pos.Longitude-l.longitude) > 1e-5
			latError := math.Abs(pos.Latitude-l.latitude) > 1e-5
			if longError || latError {
				return fmt.Errorf("Incorrect position for %q, expected (%.6f, %.6f), got (%.6f, %.6f)", l.key, l.longitude, l.latitude, pos.Longitude, pos.Latitude)
			}
			logger.Successf("Correct position for %q: (%.6f, %.6f)", l.key, pos.Longitude, pos.Latitude)
		}
	}

	return nil
}

func testGeospatialPosSingle(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	locations := random.RandomElementsFromArray(randomLocations, 1)

	tests := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities",
			locations: locations,
		},
		NewGeoPosTest("cities", locations),
	}

	return runTestCases(tests, client, logger)
}
