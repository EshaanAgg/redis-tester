package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialAddInvalid(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	testCases := []RunnableTestCase{
		&GeoAddTest{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_long", -180.000123, 12.345678},
			},
		},
		&GeoAddTest{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_long_2", 180.000123, 12.345678},
			},
		},
		&GeoAddTest{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_lat", 12.345678, -85.101234},
			},
		},
		&GeoAddTest{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_lat_2", 12.345678, 85.101234},
			},
		},
		&GeoAddTest{
			key:       "cities",
			locations: random.RandomElementsFromArray(randomLocations, 1),
		},
	}

	return runTestCases(testCases, client, logger)
}
