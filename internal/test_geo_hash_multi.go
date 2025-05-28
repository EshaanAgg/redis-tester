package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialHashMulti(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	locations1 := random.RandomElementsFromArray(randomLocations, 3)
	locations2 := random.RandomElementsFromArray(randomLocations, 5)

	testCases := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities-1",
			locations: locations1,
		},
		&GeoHashTest{
			key:       "cities-1",
			locations: random.RandomElementsFromArray(locations1, 2),
		},
		&GeoAddTest{
			key:       "cities-2",
			locations: locations2,
		},
		&GeoHashTest{
			key:       "cities-2",
			locations: random.RandomElementsFromArray(locations2, 4),
		},
	}

	return runTestCases(testCases, client, logger)
}
