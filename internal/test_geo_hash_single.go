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

type GeoHashTest struct {
	key       string
	locations []GeoLocation
}

// Returns the command string for the GeoHashTest in the format "geoadd key longitude latitude"
func (t *GeoHashTest) command() string {
	var members []string
	for _, l := range t.locations {
		members = append(members, l.key)
	}
	return fmt.Sprintf("geohash %s %s", t.key, strings.Join(members, " "))
}

func (t *GeoHashTest) Run(client *redis.Client, logger *logger.Logger) error {
	logger.Infof("$ redis-cli %s", t.command())

	var memberNames []string
	for _, l := range t.locations {
		memberNames = append(memberNames, l.key)
	}
	resp, err := client.GeoHash(t.key, memberNames...).Result()

	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	logger.Infof("Recieved hashes: %v", resp)

	// Check the response
	if len(resp) != len(t.locations) {
		return fmt.Errorf("Expected %d geohashes, got %d", len(t.locations), len(resp))
	}

	for i, l := range t.locations {
		if err := l.validateGeoHash(resp[i]); err != nil {
			logger.Errorf("Incorrect geohash for %q, expected hash to be %s, recieved %s.", l.key, GetGeoHash(l.longitude, l.latitude), resp[i])
			return err
		}
		logger.Successf("Correct geohash for %q: %q", l.key, resp[i])
	}

	return nil
}

func testGeospatialHashSingle(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	location := random.RandomElementsFromArray(randomLocations, 1)

	testCases := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities",
			locations: location,
		},
		&GeoHashTest{
			key:       "cities",
			locations: location,
		},
	}

	return runTestCases(testCases, client, logger)
}
