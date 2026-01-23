package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// RightsService Tests - TDD Red Phase
// These tests are designed to FAIL until the implementation is complete
// =============================================================================

func TestRightsService_GetTrackRights(t *testing.T) {
	tests := []struct {
		name        string
		trackID     string
		expectCount int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "existing track with rights",
			trackID:     "track-123",
			expectCount: 2, // e.g., mechanical + performance
			expectError: false,
		},
		{
			name:        "track with no rights",
			trackID:     "track-no-rights",
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "non-existent track",
			trackID:     "non-existent",
			expectError: true,
			errorMsg:    "track not found",
		},
		{
			name:        "empty track ID",
			trackID:     "",
			expectError: true,
			errorMsg:    "track ID required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil) // Will fail - service doesn't exist

			rights, err := svc.GetTrackRights(ctx, tt.trackID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Len(t, rights, tt.expectCount)
			}
		})
	}
}

func TestRightsService_SetTrackRights(t *testing.T) {
	tests := []struct {
		name        string
		trackID     string
		rights      []TrackRightsInput
		expectError bool
		errorMsg    string
	}{
		{
			name:    "valid rights with 100% share",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "mechanical", SharePercent: 50.0, Territories: []string{"US"}},
				{HolderID: "holder-2", RightType: "mechanical", SharePercent: 50.0, Territories: []string{"US"}},
			},
			expectError: false,
		},
		{
			name:    "single holder with 100%",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "performance", SharePercent: 100.0, Territories: []string{"WW"}},
			},
			expectError: false,
		},
		{
			name:    "share sum exceeds 100%",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "mechanical", SharePercent: 60.0, Territories: []string{"US"}},
				{HolderID: "holder-2", RightType: "mechanical", SharePercent: 60.0, Territories: []string{"US"}},
			},
			expectError: true,
			errorMsg:    "share percentages must sum to 100",
		},
		{
			name:    "share sum below 100%",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "mechanical", SharePercent: 40.0, Territories: []string{"US"}},
				{HolderID: "holder-2", RightType: "mechanical", SharePercent: 30.0, Territories: []string{"US"}},
			},
			expectError: true,
			errorMsg:    "share percentages must sum to 100",
		},
		{
			name:    "multiple right types - each must sum to 100%",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "mechanical", SharePercent: 100.0, Territories: []string{"US"}},
				{HolderID: "holder-1", RightType: "performance", SharePercent: 50.0, Territories: []string{"US"}},
				{HolderID: "holder-2", RightType: "performance", SharePercent: 50.0, Territories: []string{"US"}},
			},
			expectError: false,
		},
		{
			name:        "empty track ID",
			trackID:     "",
			rights:      []TrackRightsInput{},
			expectError: true,
			errorMsg:    "track ID required",
		},
		{
			name:    "invalid right type",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "invalid", SharePercent: 100.0, Territories: []string{"US"}},
			},
			expectError: true,
			errorMsg:    "invalid right type",
		},
		{
			name:    "empty territories",
			trackID: "track-123",
			rights: []TrackRightsInput{
				{HolderID: "holder-1", RightType: "mechanical", SharePercent: 100.0, Territories: []string{}},
			},
			expectError: true,
			errorMsg:    "territories required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			err := svc.SetTrackRights(ctx, tt.trackID, tt.rights)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRightsService_CheckAccess(t *testing.T) {
	tests := []struct {
		name            string
		trackID         string
		userID          string
		territory       string
		expectedAllowed bool
		expectedReason  string
		expectError     bool
	}{
		{
			name:            "allowed access in licensed territory",
			trackID:         "track-123",
			userID:          "user-456",
			territory:       "US",
			expectedAllowed: true,
			expectError:     false,
		},
		{
			name:            "geo-blocked territory",
			trackID:         "track-123",
			userID:          "user-456",
			territory:       "CN", // China - hypothetically blocked
			expectedAllowed: false,
			expectedReason:  "content unavailable in your region",
			expectError:     false,
		},
		{
			name:            "expired rights",
			trackID:         "track-expired",
			userID:          "user-456",
			territory:       "US",
			expectedAllowed: false,
			expectedReason:  "rights expired",
			expectError:     false,
		},
		{
			name:            "no rights defined - deny",
			trackID:         "track-no-rights",
			userID:          "user-456",
			territory:       "US",
			expectedAllowed: false,
			expectedReason:  "no rights information available",
			expectError:     false,
		},
		{
			name:            "worldwide rights - allow anywhere",
			trackID:         "track-worldwide",
			userID:          "user-456",
			territory:       "JP",
			expectedAllowed: true,
			expectError:     false,
		},
		{
			name:        "empty track ID",
			trackID:     "",
			userID:      "user-456",
			territory:   "US",
			expectError: true,
		},
		{
			name:        "empty territory",
			trackID:     "track-123",
			userID:      "user-456",
			territory:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			result, err := svc.CheckAccess(ctx, tt.trackID, tt.userID, tt.territory)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedAllowed, result.Allowed)
				if tt.expectedReason != "" {
					assert.Contains(t, result.Reason, tt.expectedReason)
				}
				assert.Equal(t, tt.territory, result.Territory)
			}
		})
	}
}

func TestRightsService_GetRightsHolders(t *testing.T) {
	tests := []struct {
		name        string
		trackID     string
		expectCount int
		expectError bool
	}{
		{
			name:        "track with multiple holders",
			trackID:     "track-123",
			expectCount: 3, // label, publisher, artist
			expectError: false,
		},
		{
			name:        "track with single holder",
			trackID:     "track-single",
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "track with no holders",
			trackID:     "track-no-holders",
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "empty track ID",
			trackID:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			holders, err := svc.GetRightsHolders(ctx, tt.trackID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, holders, tt.expectCount)
			}
		})
	}
}

func TestRightsService_CreateRightsHolder(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateRightsHolderInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid label holder",
			input: CreateRightsHolderInput{
				Name:        "Universal Music",
				Type:        "label",
				Territories: []string{"US", "CA", "UK"},
			},
			expectError: false,
		},
		{
			name: "valid artist holder with IPI",
			input: CreateRightsHolderInput{
				Name:        "John Doe",
				Type:        "artist",
				Territories: []string{"WW"},
				IPINumber:   strPtr("00012345678"),
			},
			expectError: false,
		},
		{
			name: "missing name",
			input: CreateRightsHolderInput{
				Type:        "label",
				Territories: []string{"US"},
			},
			expectError: true,
			errorMsg:    "name required",
		},
		{
			name: "invalid type",
			input: CreateRightsHolderInput{
				Name:        "Test Holder",
				Type:        "invalid",
				Territories: []string{"US"},
			},
			expectError: true,
			errorMsg:    "invalid holder type",
		},
		{
			name: "invalid IPI format",
			input: CreateRightsHolderInput{
				Name:        "John Doe",
				Type:        "artist",
				Territories: []string{"US"},
				IPINumber:   strPtr("12345"), // Too short
			},
			expectError: true,
			errorMsg:    "invalid IPI number format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			holder, err := svc.CreateRightsHolder(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, holder)
				assert.NotEmpty(t, holder.ID)
				assert.Equal(t, tt.input.Name, holder.Name)
			}
		})
	}
}

func TestRightsService_UpdateRightsHolder(t *testing.T) {
	tests := []struct {
		name        string
		holderID    string
		input       UpdateRightsHolderInput
		expectError bool
	}{
		{
			name:     "update name",
			holderID: "holder-123",
			input: UpdateRightsHolderInput{
				Name: strPtr("Updated Music Label"),
			},
			expectError: false,
		},
		{
			name:     "add IPI number",
			holderID: "holder-123",
			input: UpdateRightsHolderInput{
				IPINumber: strPtr("00012345678"),
			},
			expectError: false,
		},
		{
			name:     "update territories",
			holderID: "holder-123",
			input: UpdateRightsHolderInput{
				Territories: []string{"US", "CA", "UK", "DE"},
			},
			expectError: false,
		},
		{
			name:     "deactivate holder",
			holderID: "holder-123",
			input: UpdateRightsHolderInput{
				IsActive: boolPtr(false),
			},
			expectError: false,
		},
		{
			name:        "non-existent holder",
			holderID:    "non-existent",
			input:       UpdateRightsHolderInput{Name: strPtr("Test")},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			holder, err := svc.UpdateRightsHolder(ctx, tt.holderID, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, holder)
			}
		})
	}
}

func TestRightsService_TransferRights(t *testing.T) {
	tests := []struct {
		name          string
		trackID       string
		rightType     string
		fromHolderID  string
		toHolderID    string
		sharePercent  float64
		effectiveDate time.Time
		expectError   bool
	}{
		{
			name:          "transfer full rights",
			trackID:       "track-123",
			rightType:     "mechanical",
			fromHolderID:  "holder-old",
			toHolderID:    "holder-new",
			sharePercent:  100.0,
			effectiveDate: time.Now().Add(30 * 24 * time.Hour),
			expectError:   false,
		},
		{
			name:          "transfer partial rights",
			trackID:       "track-123",
			rightType:     "mechanical",
			fromHolderID:  "holder-old",
			toHolderID:    "holder-new",
			sharePercent:  25.0,
			effectiveDate: time.Now().Add(30 * 24 * time.Hour),
			expectError:   false,
		},
		{
			name:          "transfer more than owned",
			trackID:       "track-123",
			rightType:     "mechanical",
			fromHolderID:  "holder-partial", // owns 50%
			toHolderID:    "holder-new",
			sharePercent:  75.0,
			effectiveDate: time.Now(),
			expectError:   true,
		},
		{
			name:          "transfer to self",
			trackID:       "track-123",
			rightType:     "mechanical",
			fromHolderID:  "holder-1",
			toHolderID:    "holder-1",
			sharePercent:  50.0,
			effectiveDate: time.Now(),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewRightsService(nil, nil)

			err := svc.TransferRights(ctx, tt.trackID, tt.rightType, tt.fromHolderID, tt.toHolderID, tt.sharePercent, tt.effectiveDate)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions for tests
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// Input types that should be defined in rights.go
type TrackRightsInput struct {
	HolderID     string   `json:"holderId"`
	RightType    string   `json:"rightType"`
	SharePercent float64  `json:"sharePercent"`
	Territories  []string `json:"territories"`
	StartDate    *time.Time `json:"startDate,omitempty"`
	EndDate      *time.Time `json:"endDate,omitempty"`
}

type CreateRightsHolderInput struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Territories []string `json:"territories"`
	IPINumber   *string  `json:"ipiNumber,omitempty"`
	ISNI        *string  `json:"isni,omitempty"`
	ArtistID    *string  `json:"artistId,omitempty"`
}

type UpdateRightsHolderInput struct {
	Name        *string  `json:"name,omitempty"`
	Territories []string `json:"territories,omitempty"`
	IPINumber   *string  `json:"ipiNumber,omitempty"`
	ISNI        *string  `json:"isni,omitempty"`
	IsActive    *bool    `json:"isActive,omitempty"`
}
