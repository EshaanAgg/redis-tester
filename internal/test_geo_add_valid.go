package internal

import (
	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testGeospatialGeoaddValid(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	testCases := []GeoAddTest{
		{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_long", -180.000123, 12.345678},
			},
			expectedError: "ERR invalid longitude,latitude pair -180.000123,12.345678",
		},
		{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_long_2", 180.000123, 12.345678},
			},
			expectedError: "ERR invalid longitude,latitude pair 180.000123,12.345678",
		},
		{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_lat", 12.345678, -85.101234},
			},
			expectedError: "ERR invalid longitude,latitude pair 12.345678,-85.101234",
		},
		{
			key: "cities",
			locations: []GeoLocation{
				{"invalid_lat_2", 12.345678, 85.101234},
			},
			expectedError: "ERR invalid longitude,latitude pair 12.345678,85.101234",
		},
		{
			key:       "cities",
			locations: random.RandomElementsFromArray(randomLocations, 1),
		},
	}

	for _, tc := range testCases {
		if err := tc.Run(client, logger); err != nil {
			return err
		}
	}

	return nil
}
