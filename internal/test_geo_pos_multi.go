package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialPosMulti(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	locations := random.RandomElementsFromArray(randomLocations, 4)
	queryOne := random.RandomElementsFromArray(locations, 1)
	queryThree := random.RandomElementsFromArray(locations, 3)
	queryWithInvalid := append(
		getGeoPositionsFromLocations(random.RandomElementsFromArray(locations, 2)),
		GeoPosition{key: "invalid_city", isInvalid: true},
	)

	tests := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities",
			locations: locations,
		},
		NewGeoPosTest("cities", queryOne),
		&GeoPosTest{
			key:       "cities",
			locations: queryWithInvalid,
		},
		NewGeoPosTest("cities", queryThree),
	}

	return runTestCases(tests, client, logger)
}
