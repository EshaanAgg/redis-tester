package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialGeoaddMulti(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	testCases := []GeoAddTest{
		{
			key:       "cities_1",
			locations: random.RandomElementsFromArray(randomLocations, 2),
		},
		{
			key:       "cities_2",
			locations: random.RandomElementsFromArray(randomLocations, 5),
		},
		// Should reject if any one of the locations is invalid
		{
			key: "invalid_city",
			locations: []GeoLocation{
				random.RandomElementFromArray(randomLocations),
				{"invalid_lat", 91.123456, 181.123456},
			},
			expectedError: "ERR invalid longitude,latitude pair 91.123456,181.123456",
		},
	}

	for _, tc := range testCases {
		if err := tc.Run(client, logger); err != nil {
			return err
		}
	}

	return nil
}
