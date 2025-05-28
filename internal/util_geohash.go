package internal

import (
	"fmt"
	"strings"
)

const GEOHASH_LENGTH = 11
const MAXIMUM_GEO_HASH_BITS = 26

// getGeoHashBits recursively calculates the bits for a geohash based on the value's position
// within the specified range (lower to upper).
func getGeoHashBits(lower, upper, value float64, currentBits int64, bitCnt int64) int64 {
	if bitCnt >= MAXIMUM_GEO_HASH_BITS {
		return currentBits
	}

	mid := (lower + upper) / 2
	if value < mid {
		currentBits = (currentBits << 1) // Add 0
		return getGeoHashBits(lower, mid, value, currentBits, bitCnt+1)
	} else {
		currentBits = (currentBits << 1) | 1 // Add 1
		return getGeoHashBits(mid, upper, value, currentBits, bitCnt+1)
	}
}

func encodeGeoHash(geoBits int64) string {
	base32Chars := "0123456789bcdefghjkmnpqrstuvwxyz"
	var encoded string

	totalBits := MAXIMUM_GEO_HASH_BITS * 2         // 52 bits total
	unusedBits := (GEOHASH_LENGTH * 5) - totalBits // 3 unused bits at the bottom

	for i := 0; i < GEOHASH_LENGTH; i++ {
		shift := max(totalBits-((i+1)*5)+unusedBits, 0)
		index := (geoBits >> shift) & 0x1F
		encoded += string(base32Chars[index])
	}
	return encoded
}

func GetGeoHash(longitude, latitude float64) string {
	longitudeBits := getGeoHashBits(-180, 180, longitude, 0, 0)
	latitudeBits := getGeoHashBits(-90, 90, latitude, 0, 0)

	// Combine the bits into a single geohash by interleaving
	// Standard geohash format: longitude bit first, then latitude bit
	geoBits := int64(0)
	for bitIdx := 0; bitIdx < MAXIMUM_GEO_HASH_BITS; bitIdx++ {
		// Get bits from most significant to least significant
		bitLongitude := (longitudeBits >> (MAXIMUM_GEO_HASH_BITS - 1 - bitIdx)) & 1
		bitLatitude := (latitudeBits >> (MAXIMUM_GEO_HASH_BITS - 1 - bitIdx)) & 1

		// Interleave: longitude bit first (more significant), then latitude bit
		geoBits = (geoBits << 1) | bitLongitude
		geoBits = (geoBits << 1) | bitLatitude
	}

	return encodeGeoHash(geoBits)
}

func DecodeGeoHash(geohash string) (float64, float64, error) {
	if len(geohash) != GEOHASH_LENGTH {
		return 0, 0, fmt.Errorf("geohash must be %d characters long, but got %d characters", GEOHASH_LENGTH, len(geohash))
	}

	// Decode the geohash into bits
	base32Chars := "0123456789bcdefghjkmnpqrstuvwxyz"
	geoBits := int64(0)
	for i, char := range geohash {
		index := strings.IndexByte(base32Chars, byte(char))
		if index < 0 {
			return 0, 0, fmt.Errorf("invalid character '%c' in geohash at index %d", char, i)
		}
		geoBits = (geoBits << 5) | int64(index)
	}

	// Extract longitude and latitude bits from the interleaved geohash bits
	longitudeBits := int64(0)
	latitudeBits := int64(0)

	totalBits := MAXIMUM_GEO_HASH_BITS * 2         // 52 bits total
	unusedBits := (GEOHASH_LENGTH * 5) - totalBits // 3 unused bits at the bottom

	// Remove unused bits
	geoBits = geoBits >> unusedBits

	// De-interleave the bits
	for i := 0; i < MAXIMUM_GEO_HASH_BITS; i++ {
		// Extract pairs of bits (longitude first, then latitude)
		bitIdx := (MAXIMUM_GEO_HASH_BITS - 1 - i) * 2
		longitudeBit := (geoBits >> (bitIdx + 1)) & 1
		latitudeBit := (geoBits >> bitIdx) & 1

		longitudeBits = (longitudeBits << 1) | longitudeBit
		latitudeBits = (latitudeBits << 1) | latitudeBit
	}

	// Convert bits back to coordinates
	longitude := -180.0 + (float64(longitudeBits)*360.0)/float64(int64(1)<<MAXIMUM_GEO_HASH_BITS)
	latitude := -90.0 + (float64(latitudeBits)*180.0)/float64(int64(1)<<MAXIMUM_GEO_HASH_BITS)

	return longitude, latitude, nil
}
