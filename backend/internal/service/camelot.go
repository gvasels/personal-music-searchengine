package service

// CamelotWheel represents the Camelot key compatibility system for harmonic mixing.
// Each key maps to harmonically compatible keys (same key, adjacent on wheel, relative major/minor).
var CamelotWheel = map[string][]string{
	// Minor keys (A)
	"1A":  {"1A", "12A", "2A", "1B"},
	"2A":  {"2A", "1A", "3A", "2B"},
	"3A":  {"3A", "2A", "4A", "3B"},
	"4A":  {"4A", "3A", "5A", "4B"},
	"5A":  {"5A", "4A", "6A", "5B"},
	"6A":  {"6A", "5A", "7A", "6B"},
	"7A":  {"7A", "6A", "8A", "7B"},
	"8A":  {"8A", "7A", "9A", "8B"},
	"9A":  {"9A", "8A", "10A", "9B"},
	"10A": {"10A", "9A", "11A", "10B"},
	"11A": {"11A", "10A", "12A", "11B"},
	"12A": {"12A", "11A", "1A", "12B"},
	// Major keys (B)
	"1B":  {"1B", "12B", "2B", "1A"},
	"2B":  {"2B", "1B", "3B", "2A"},
	"3B":  {"3B", "2B", "4B", "3A"},
	"4B":  {"4B", "3B", "5B", "4A"},
	"5B":  {"5B", "4B", "6B", "5A"},
	"6B":  {"6B", "5B", "7B", "6A"},
	"7B":  {"7B", "6B", "8B", "7A"},
	"8B":  {"8B", "7B", "9B", "8A"},
	"9B":  {"9B", "8B", "10B", "9A"},
	"10B": {"10B", "9B", "11B", "10A"},
	"11B": {"11B", "10B", "12B", "11A"},
	"12B": {"12B", "11B", "1B", "12A"},
}

// KeyTransitions describes the type of energy change when mixing from one key to another.
var KeyTransitions = map[string]string{
	"same":       "Perfect Match",
	"adjacent":   "Smooth Transition",
	"relative":   "Major/Minor Switch",
	"compatible": "Harmonic Mix",
}

// IsKeyCompatible checks if two Camelot keys can be mixed harmonically.
// Returns true if key2 is in the compatible keys list for key1.
func IsKeyCompatible(key1, key2 string) bool {
	if key1 == "" || key2 == "" {
		return false
	}
	compatibleKeys, ok := CamelotWheel[key1]
	if !ok {
		return false
	}
	for _, k := range compatibleKeys {
		if k == key2 {
			return true
		}
	}
	return false
}

// GetCompatibleKeys returns all keys that can be harmonically mixed with the given key.
// Returns nil if the key is invalid.
func GetCompatibleKeys(key string) []string {
	if key == "" {
		return nil
	}
	compatibleKeys, ok := CamelotWheel[key]
	if !ok {
		return nil
	}
	// Return a copy to prevent modification
	result := make([]string, len(compatibleKeys))
	copy(result, compatibleKeys)
	return result
}

// GetKeyTransition describes the type of transition when mixing from one key to another.
// Returns empty string if keys are not compatible.
func GetKeyTransition(fromKey, toKey string) string {
	if fromKey == "" || toKey == "" {
		return ""
	}
	if !IsKeyCompatible(fromKey, toKey) {
		return ""
	}

	// Same key
	if fromKey == toKey {
		return KeyTransitions["same"]
	}

	// Check for relative major/minor (same number, different letter)
	fromNum := fromKey[:len(fromKey)-1]
	fromLetter := fromKey[len(fromKey)-1:]
	toNum := toKey[:len(toKey)-1]
	toLetter := toKey[len(toKey)-1:]

	if fromNum == toNum && fromLetter != toLetter {
		return KeyTransitions["relative"]
	}

	// Must be adjacent
	return KeyTransitions["adjacent"]
}

// GetBPMCompatibility checks if two BPM values are compatible for mixing.
// Returns the absolute difference and whether they're within the tolerance.
func GetBPMCompatibility(bpm1, bpm2, tolerance int) (diff int, compatible bool) {
	if bpm1 <= 0 || bpm2 <= 0 {
		return 0, false
	}

	diff = bpm1 - bpm2
	if diff < 0 {
		diff = -diff
	}

	// Also check half/double time compatibility
	halfTime1 := bpm1 / 2
	doubleTime1 := bpm1 * 2
	halfTime2 := bpm2 / 2
	doubleTime2 := bpm2 * 2

	// Direct BPM match
	if diff <= tolerance {
		return diff, true
	}

	// Half time match (e.g., 140 BPM with 70 BPM)
	halfDiff := halfTime1 - bpm2
	if halfDiff < 0 {
		halfDiff = -halfDiff
	}
	if halfDiff <= tolerance {
		return halfDiff, true
	}

	halfDiff = bpm1 - halfTime2
	if halfDiff < 0 {
		halfDiff = -halfDiff
	}
	if halfDiff <= tolerance {
		return halfDiff, true
	}

	// Double time match (rarely used but valid)
	doubleDiff := doubleTime1 - bpm2
	if doubleDiff < 0 {
		doubleDiff = -doubleDiff
	}
	if doubleDiff <= tolerance {
		return doubleDiff, true
	}

	doubleDiff = bpm1 - doubleTime2
	if doubleDiff < 0 {
		doubleDiff = -doubleDiff
	}
	if doubleDiff <= tolerance {
		return doubleDiff, true
	}

	return diff, false
}
