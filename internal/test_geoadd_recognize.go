package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/go-redis/redis"
)

type GeoLocation struct {
	key       string
	longitude float64
	latitude  float64
}

var randomLocations []GeoLocation = []GeoLocation{
	{"London", -0.1278, 51.5074},
	{"Tokyo", 139.6917, 35.6895},
	{"Paris", 2.3522, 48.8566},
	{"Mumbai", 72.8777, 19.0760},
	{"Shanghai", 121.4854, 31.2351},
	{"Beijing", 116.4074, 39.9042},
	{"Istanbul", 28.9784, 41.0082},
}

// Returns the command string for the GeoLocation in the format "longitude latitude key"
// Latitude and longitude are specified upto 4 decimal places.
func (m *GeoLocation) command() string {
	return fmt.Sprintf("%.4f %.4f %s", m.longitude, m.latitude, m.key)
}

func (m *GeoLocation) toRedisGeoLocation() *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      m.key,
		Longitude: m.longitude,
		Latitude:  m.latitude,
	}
}

type GeoAddTest struct {
	key           string
	locations     []GeoLocation
	expectedError string
}

func (t *GeoAddTest) command() string {
	var commands []string
	for _, l := range t.locations {
		commands = append(commands, l.command())
	}
	return fmt.Sprintf("geoadd %s %s", t.key, strings.Join(commands, " "))
}

func (t *GeoAddTest) Run(client *redis.Client, logger *logger.Logger) error {
	logger.Infof("$ redis-cli %s", t.command())

	var locations []*redis.GeoLocation
	for _, l := range t.locations {
		locations = append(locations, l.toRedisGeoLocation())
	}
	resp, err := client.GeoAdd(t.key, locations...).Result()

	// Handle the errors
	if err != nil && t.expectedError == "" {
		logFriendlyError(logger, err)
		return err
	}

	if err != nil && t.expectedError != "" {
		if err.Error() != t.expectedError {
			return fmt.Errorf("Expected %q, got %q", t.expectedError, err.Error())
		}

		logger.Successf("Received error: %q", err.Error())
		return nil
	}

	expectedResponse := int64(len(locations))
	if resp != expectedResponse && t.expectedError != "" {
		logger.Infof("Received response: %q", resp)
		return fmt.Errorf("Expected an error as the response, got %q", resp)
	} else if resp != expectedResponse {
		logger.Infof("Received response: %q", resp)
		return fmt.Errorf("Expected %q, got %q", expectedResponse, resp)
	} else {
		logger.Successf("Received response: %q", resp)
	}

	return nil
}

func testGeospatialGeoaddRecognize(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	tc := GeoAddTest{
		key:       "cities",
		locations: random.RandomElementsFromArray(randomLocations, 1),
	}
	return tc.Run(client, logger)
}
