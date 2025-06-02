package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialDisUnits(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	locations := random.RandomElementsFromArray(randomLocations, 5)

	tests := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities",
			locations: locations,
		},
	}

	for _, unit := range allGeoDistUnits {
		randomPoints := random.RandomElementsFromArray(locations, 2)
		tests = append(tests, &GeoDisTest{
			key:     "cities",
			memberA: randomPoints[0],
			memberB: randomPoints[1],
			unit:    unit,
		})
	}

	return runTestCases(tests, client, logger)
}
