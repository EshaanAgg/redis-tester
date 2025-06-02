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

const MAXIMUM_GEOHASH_ERROR float64 = 1e-5

type GeoLocation struct {
	key       string
	longitude float64
	latitude  float64
}

var randomLocations []GeoLocation = []GeoLocation{
	{"london", -0.127834, 51.507456},
	{"tokyo", 139.691754, 35.689512},
	{"paris", 2.352224, 48.856689},
	{"mumbai", 72.877727, 19.076032},
	{"shanghai", 121.485487, 31.235182},
	{"beijing", 116.407449, 39.904227},
	{"istanbul", 28.978447, 41.008235},
}

// Returns the command string for the GeoLocation in the format "longitude latitude key"
// Latitude and longitude are specified upto 6 decimal places.
func (m *GeoLocation) command() string {
	return fmt.Sprintf("%.6f %.6f %s", m.longitude, m.latitude, m.key)
}

func (m *GeoLocation) validateGeoHash(hash string) error {
	long, lat, err := DecodeGeoHash(hash)
	if err != nil {
		return fmt.Errorf("The geohash '%s' could not be decoded: %v", hash, err)
	}

	errLong := math.Abs(long - m.longitude)
	errLat := math.Abs(lat - m.latitude)

	if math.Max(errLong, errLat) > MAXIMUM_GEOHASH_ERROR {
		errMsg := fmt.Sprintf("The decoded location (%.6f, %.6f) does not match the original location (%.6f, %.6f) for geohash '%s'",
			long, lat, m.longitude, m.latitude, hash)

		if errLong > MAXIMUM_GEOHASH_ERROR {
			errMsg += fmt.Sprintf("\nLongitude error: %.6f", errLong)
		}
		if errLat > MAXIMUM_GEOHASH_ERROR {
			errMsg += fmt.Sprintf("\nLatitude error: %.6f", errLat)
		}

		if math.Max(errLong, errLat) < 1e-3 {
			errMsg += "\nThis might be due to rounding errors in the geohash encoding/decoding process. Please ensure that you are using all 52 bits of the geohash during the encoding and decoding process."
		}

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (m *GeoLocation) toRedisGeoLocation() *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      m.key,
		Longitude: m.longitude,
		Latitude:  m.latitude,
	}
}

type GeoAddTest struct {
	key       string
	locations []GeoLocation

	expectedError string
}

func (t *GeoAddTest) command() string {
	var commands []string
	for _, l := range t.locations {
		commands = append(commands, l.command())
	}
	return fmt.Sprintf("geoadd %s %s", t.key, strings.Join(commands, " "))
}

// Validates the GeoAddTest by checking if the locations are within valid ranges.
// If any location is invalid, it sets the expectedError field with a descriptive error message.
func (t *GeoAddTest) validate() {
	for _, l := range t.locations {
		if l.longitude < -180 || l.longitude > 180 || l.latitude < -85.05112878 || l.latitude > 85.05112878 {
			t.expectedError = fmt.Sprintf("ERR invalid longitude,latitude pair %.6f,%.6f", l.longitude, l.latitude)
			return
		}
	}
}

func (t *GeoAddTest) Run(client *redis.Client, logger *logger.Logger) error {
	t.validate()
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

	// Check the response
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

func testGeospatialAddSingle(stageHarness *test_case_harness.TestCaseHarness) error {
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
