package internal

import (
	"fmt"
	"math"

	"github.com/codecrafters-io/redis-tester/internal/redis_executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/go-redis/redis"
)

type GeoDistUnit string

const (
	GeoDistUnitMeters     GeoDistUnit = "m"
	GeoDistUnitKilometers GeoDistUnit = "km"
	GeoDistUnitMiles      GeoDistUnit = "mi"
	GeoDistUnitFeet       GeoDistUnit = "ft"
)

var allGeoDistUnits = []GeoDistUnit{
	GeoDistUnitMeters,
	GeoDistUnitKilometers,
	GeoDistUnitMiles,
	GeoDistUnitFeet,
}

func (d *GeoDistUnit) String() string {
	switch *d {
	case GeoDistUnitMeters:
		return "m"
	case GeoDistUnitKilometers:
		return "km"
	case GeoDistUnitMiles:
		return "mi"
	case GeoDistUnitFeet:
		return "ft"
	default:
		panic(fmt.Sprintf("Unknown GeoDistUnit: %s", *d))
	}
}

func (d *GeoDistUnit) getConvertFactor() float64 {
	switch *d {
	case GeoDistUnitMeters:
		return 1.0
	case GeoDistUnitKilometers:
		return 1000.0
	case GeoDistUnitMiles:
		return 1609.34
	case GeoDistUnitFeet:
		return 0.3048
	default:
		panic(fmt.Sprintf("Unknown GeoDistUnit: %s", *d))
	}
}

type GeoDisTest struct {
	key     string
	memberA GeoLocation
	memberB GeoLocation
	unit    GeoDistUnit
}

func (g *GeoDisTest) command() string {
	if g.unit == GeoDistUnitMeters {
		return fmt.Sprintf("geodist %s %s %s", g.key, g.memberA.key, g.memberB.key)
	}
	return fmt.Sprintf("geodist %s %s %s %s", g.key, g.memberA.key, g.memberB.key, g.unit)
}

// Returns the distance between two members assuming the Earth to be a perfect sphere.
// The formula used is based on the Haversine formula.
func (g *GeoDisTest) getDistance() float64 {
	// Convert degrees to radians
	latA := g.memberA.latitude * math.Pi / 180.0
	lonA := g.memberA.longitude * math.Pi / 180.0
	latB := g.memberB.latitude * math.Pi / 180.0
	lonB := g.memberB.longitude * math.Pi / 180.0

	// Haversine formula
	dLat := latB - latA
	dLon := lonB - lonA

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(latA)*math.Cos(latB)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Radius of the Earth in meters
	radius := 6371000.0
	distM := radius * c                      // Distance in meters
	return distM / g.unit.getConvertFactor() // Convert to the requested unit
}

func (g *GeoDisTest) Run(client *redis.Client, logger *logger.Logger) error {
	logger.Infof("$ redis-cli %s", g.command())

	resp, err := client.GeoDist(g.key, g.memberA.key, g.memberB.key, g.unit.String()).Result()
	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	// Allow upto 1% error margin for floating point calculations
	const errorMargin = 0.01
	expectedDistance := g.getDistance()
	difference := math.Abs(resp - expectedDistance)
	if difference > expectedDistance*errorMargin {
		errMsg := fmt.Sprintf("Expected distance: %.6f %s, got: %.6f %s (difference: %.6f %s)",
			expectedDistance, g.unit, resp, g.unit, difference, g.unit)
		return fmt.Errorf(errMsg)
	}

	errPercentage := fmt.Sprintf("%.2f%%", (difference/expectedDistance)*100)
	logger.Successf("Correctly logged distance: %.6f %s (%s error)", resp, g.unit, errPercentage)

	return nil
}

func testGeospatialDis(stageHarness *test_case_harness.TestCaseHarness) error {
	b := redis_executable.NewRedisExecutable(stageHarness)
	if err := b.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger
	client := NewRedisClient("localhost:6379")
	defer client.Close()

	locations := random.RandomElementsFromArray(randomLocations, 2)

	tests := []RunnableTestCase{
		&GeoAddTest{
			key:       "cities",
			locations: locations,
		},
		&GeoDisTest{
			key:     "cities",
			memberA: locations[0],
			memberB: locations[1],
			unit:    GeoDistUnitMeters, // Default unit is meters
		},
	}

	return runTestCases(tests, client, logger)
}
