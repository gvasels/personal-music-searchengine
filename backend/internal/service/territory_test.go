package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TerritoryService Tests - TDD Red Phase
// These tests are designed to FAIL until the implementation is complete
// =============================================================================

func TestTerritoryService_GetTerritory(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid US code",
			code:        "US",
			expectError: false,
		},
		{
			name:        "valid regional code US-CA",
			code:        "US-CA",
			expectError: false,
		},
		{
			name:        "valid venue code",
			code:        "LOC:VENUE123",
			expectError: false,
		},
		{
			name:        "valid worldwide code",
			code:        "WW",
			expectError: false,
		},
		{
			name:        "invalid code",
			code:        "INVALID",
			expectError: true,
			errorMsg:    "territory not found",
		},
		{
			name:        "empty code",
			code:        "",
			expectError: true,
			errorMsg:    "territory code required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil) // Will fail - service doesn't exist

			territory, err := svc.GetTerritory(ctx, tt.code)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, territory)
				assert.Equal(t, tt.code, territory.Code)
			}
		})
	}
}

func TestTerritoryService_ResolveHierarchy(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		expectedLevels int // Expected number of territories in hierarchy
		expectError    bool
	}{
		{
			name:           "US resolves to WW",
			code:           "US",
			expectedLevels: 2, // US → WW
			expectError:    false,
		},
		{
			name:           "US-CA resolves through US to WW",
			code:           "US-CA",
			expectedLevels: 3, // US-CA → US → WW
			expectError:    false,
		},
		{
			name:           "venue resolves full hierarchy",
			code:           "LOC:VENUE123",
			expectedLevels: 4, // LOC → US-CA → US → WW (example)
			expectError:    false,
		},
		{
			name:           "WW is already global",
			code:           "WW",
			expectedLevels: 1, // Just WW
			expectError:    false,
		},
		{
			name:        "invalid code",
			code:        "INVALID",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil)

			hierarchy, err := svc.ResolveHierarchy(ctx, tt.code)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, hierarchy, tt.expectedLevels)
				// First item should be the requested territory
				assert.Equal(t, tt.code, hierarchy[0].Code)
				// Last item should be global (WW)
				assert.Equal(t, "WW", hierarchy[len(hierarchy)-1].Code)
			}
		})
	}
}

func TestTerritoryService_GetPROForTerritory(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectedPRO string
		expectError bool
	}{
		{
			name:        "US has ASCAP as default PRO",
			code:        "US",
			expectedPRO: "ASCAP",
			expectError: false,
		},
		{
			name:        "UK has PRS",
			code:        "UK",
			expectedPRO: "PRS",
			expectError: false,
		},
		{
			name:        "DE has GEMA",
			code:        "DE",
			expectedPRO: "GEMA",
			expectError: false,
		},
		{
			name:        "FR has SACEM",
			code:        "FR",
			expectedPRO: "SACEM",
			expectError: false,
		},
		{
			name:        "territory with no PRO",
			code:        "XX",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil)

			pro, err := svc.GetPROForTerritory(ctx, tt.code)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, pro)
				assert.Equal(t, tt.expectedPRO, pro.Name)
			}
		})
	}
}

func TestTerritoryService_GetRoyaltyRates(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		rightType    string
		expectedRate float64
		expectError  bool
	}{
		{
			name:         "US mechanical rate",
			code:         "US",
			rightType:    "mechanical",
			expectedRate: 0.091, // 9.1 cents per stream
			expectError:  false,
		},
		{
			name:         "US performance rate",
			code:         "US",
			rightType:    "performance",
			expectedRate: 0.12, // 12% of revenue
			expectError:  false,
		},
		{
			name:         "DE sync rate (negotiated)",
			code:         "DE",
			rightType:    "sync",
			expectedRate: 0.0, // Sync is negotiated per use
			expectError:  false,
		},
		{
			name:        "invalid territory",
			code:        "INVALID",
			rightType:   "mechanical",
			expectError: true,
		},
		{
			name:        "invalid right type",
			code:        "US",
			rightType:   "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil)

			rate, err := svc.GetRoyaltyRates(ctx, tt.code, tt.rightType)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, rate)
				assert.Equal(t, tt.expectedRate, rate.Rate)
			}
		})
	}
}

func TestTerritoryService_ListTerritories(t *testing.T) {
	tests := []struct {
		name        string
		scope       string
		expectCount int
		expectError bool
	}{
		{
			name:        "list all national territories",
			scope:       "national",
			expectCount: 195, // Approximately all countries
			expectError: false,
		},
		{
			name:        "list all regional territories",
			scope:       "regional",
			expectCount: 50, // US states as example
			expectError: false,
		},
		{
			name:        "list global territory",
			scope:       "global",
			expectCount: 1, // Just WW
			expectError: false,
		},
		{
			name:        "invalid scope",
			scope:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil)

			territories, err := svc.ListTerritories(ctx, tt.scope)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(territories), 1)
			}
		})
	}
}

func TestTerritoryService_IsSubTerritory(t *testing.T) {
	tests := []struct {
		name     string
		child    string
		parent   string
		expected bool
	}{
		{
			name:     "US-CA is sub of US",
			child:    "US-CA",
			parent:   "US",
			expected: true,
		},
		{
			name:     "US is sub of WW",
			child:    "US",
			parent:   "WW",
			expected: true,
		},
		{
			name:     "US-CA is sub of WW (transitive)",
			child:    "US-CA",
			parent:   "WW",
			expected: true,
		},
		{
			name:     "US is not sub of UK",
			child:    "US",
			parent:   "UK",
			expected: false,
		},
		{
			name:     "WW is not sub of anything",
			child:    "WW",
			parent:   "US",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewTerritoryService(nil)

			result, err := svc.IsSubTerritory(ctx, tt.child, tt.parent)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
