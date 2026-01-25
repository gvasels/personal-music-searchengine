package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsKeyCompatible_SameKey(t *testing.T) {
	testCases := []string{"1A", "5A", "8A", "12A", "1B", "6B", "11B"}
	for _, key := range testCases {
		assert.True(t, IsKeyCompatible(key, key), "same key %s should be compatible", key)
	}
}

func TestIsKeyCompatible_AdjacentKeys(t *testing.T) {
	testCases := []struct {
		key1, key2 string
	}{
		{"8A", "7A"},
		{"8A", "9A"},
		{"1A", "12A"},
		{"1A", "2A"},
		{"8B", "7B"},
		{"8B", "9B"},
		{"1B", "12B"},
		{"1B", "2B"},
	}
	for _, tc := range testCases {
		assert.True(t, IsKeyCompatible(tc.key1, tc.key2),
			"%s and %s should be compatible (adjacent)", tc.key1, tc.key2)
	}
}

func TestIsKeyCompatible_RelativeMajorMinor(t *testing.T) {
	testCases := []struct {
		key1, key2 string
	}{
		{"8A", "8B"},
		{"8B", "8A"},
		{"1A", "1B"},
		{"12A", "12B"},
	}
	for _, tc := range testCases {
		assert.True(t, IsKeyCompatible(tc.key1, tc.key2),
			"%s and %s should be compatible (relative major/minor)", tc.key1, tc.key2)
	}
}

func TestIsKeyCompatible_NotCompatible(t *testing.T) {
	testCases := []struct {
		key1, key2 string
	}{
		{"8A", "3A"},
		{"8A", "4B"},
		{"1A", "6A"},
		{"5B", "10B"},
	}
	for _, tc := range testCases {
		assert.False(t, IsKeyCompatible(tc.key1, tc.key2),
			"%s and %s should NOT be compatible", tc.key1, tc.key2)
	}
}

func TestIsKeyCompatible_EmptyKeys(t *testing.T) {
	assert.False(t, IsKeyCompatible("", "8A"), "empty first key should not be compatible")
	assert.False(t, IsKeyCompatible("8A", ""), "empty second key should not be compatible")
	assert.False(t, IsKeyCompatible("", ""), "both empty keys should not be compatible")
}

func TestIsKeyCompatible_InvalidKeys(t *testing.T) {
	assert.False(t, IsKeyCompatible("invalid", "8A"), "invalid first key should not be compatible")
	assert.False(t, IsKeyCompatible("8A", "invalid"), "invalid second key should not be compatible")
	assert.False(t, IsKeyCompatible("13A", "8A"), "non-existent key should not be compatible")
}

func TestGetCompatibleKeys_AllKeys(t *testing.T) {
	// All 24 keys should have exactly 4 compatible keys
	allKeys := []string{
		"1A", "2A", "3A", "4A", "5A", "6A", "7A", "8A", "9A", "10A", "11A", "12A",
		"1B", "2B", "3B", "4B", "5B", "6B", "7B", "8B", "9B", "10B", "11B", "12B",
	}
	for _, key := range allKeys {
		compatibles := GetCompatibleKeys(key)
		assert.Len(t, compatibles, 4, "key %s should have 4 compatible keys", key)
		assert.Contains(t, compatibles, key, "compatible keys should include the key itself")
	}
}

func TestGetCompatibleKeys_EmptyKey(t *testing.T) {
	assert.Nil(t, GetCompatibleKeys(""), "empty key should return nil")
}

func TestGetCompatibleKeys_InvalidKey(t *testing.T) {
	assert.Nil(t, GetCompatibleKeys("invalid"), "invalid key should return nil")
}

func TestGetKeyTransition_SameKey(t *testing.T) {
	transition := GetKeyTransition("8A", "8A")
	assert.Equal(t, "Perfect Match", transition)
}

func TestGetKeyTransition_RelativeMajorMinor(t *testing.T) {
	transition := GetKeyTransition("8A", "8B")
	assert.Equal(t, "Major/Minor Switch", transition)

	transition = GetKeyTransition("5B", "5A")
	assert.Equal(t, "Major/Minor Switch", transition)
}

func TestGetKeyTransition_Adjacent(t *testing.T) {
	transition := GetKeyTransition("8A", "7A")
	assert.Equal(t, "Smooth Transition", transition)

	transition = GetKeyTransition("8A", "9A")
	assert.Equal(t, "Smooth Transition", transition)
}

func TestGetKeyTransition_NotCompatible(t *testing.T) {
	transition := GetKeyTransition("8A", "3A")
	assert.Empty(t, transition, "incompatible keys should return empty string")
}

func TestGetKeyTransition_EmptyKeys(t *testing.T) {
	assert.Empty(t, GetKeyTransition("", "8A"), "empty first key should return empty string")
	assert.Empty(t, GetKeyTransition("8A", ""), "empty second key should return empty string")
}

func TestGetBPMCompatibility_Direct(t *testing.T) {
	diff, compatible := GetBPMCompatibility(128, 130, 5)
	assert.True(t, compatible, "128 and 130 should be compatible with tolerance 5")
	assert.Equal(t, 2, diff)
}

func TestGetBPMCompatibility_ExactMatch(t *testing.T) {
	diff, compatible := GetBPMCompatibility(128, 128, 5)
	assert.True(t, compatible, "exact same BPM should be compatible")
	assert.Equal(t, 0, diff)
}

func TestGetBPMCompatibility_AtTolerance(t *testing.T) {
	diff, compatible := GetBPMCompatibility(128, 133, 5)
	assert.True(t, compatible, "128 and 133 should be compatible with tolerance 5")
	assert.Equal(t, 5, diff)
}

func TestGetBPMCompatibility_OutOfTolerance(t *testing.T) {
	diff, compatible := GetBPMCompatibility(128, 140, 5)
	assert.False(t, compatible, "128 and 140 should NOT be compatible with tolerance 5")
	assert.Equal(t, 12, diff)
}

func TestGetBPMCompatibility_HalfTime(t *testing.T) {
	_, compatible := GetBPMCompatibility(140, 70, 5)
	assert.True(t, compatible, "140 BPM should be compatible with 70 BPM (half time)")
}

func TestGetBPMCompatibility_ZeroBPM(t *testing.T) {
	_, compatible := GetBPMCompatibility(0, 128, 5)
	assert.False(t, compatible, "zero BPM should not be compatible")

	_, compatible = GetBPMCompatibility(128, 0, 5)
	assert.False(t, compatible, "zero BPM should not be compatible")
}
