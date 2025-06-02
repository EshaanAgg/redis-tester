package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialPosInvalid(stageHarness *test_case_harness.TestCaseHarness) error {
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
		&GeoPosTest{
			key: "cities",
			locations: []GeoPosition{
				{key: "invalid_city", isInvalid: true}, // Invalid coordinates
			},
		},
		NewGeoPosTest("cities", locations),
	}

	return runTestCases(tests, client, logger)
}
