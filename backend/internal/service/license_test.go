package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// LicenseService Tests - TDD Red Phase
// These tests are designed to FAIL until the implementation is complete
// =============================================================================

func TestLicenseService_CreateLicense(t *testing.T) {
	now := time.Now()
	future := now.Add(365 * 24 * time.Hour)

	tests := []struct {
		name        string
		input       CreateLicenseInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid sync license",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{"US", "CA"},
				StartDate:   now,
				EndDate:     future,
				AutoRenew:   false,
				Fee:         50000, // $500.00
				Currency:    "USD",
			},
			expectError: false,
		},
		{
			name: "valid mechanical license with auto-renew",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "mechanical",
				Territories: []string{"WW"},
				StartDate:   now,
				EndDate:     future,
				AutoRenew:   true,
				Fee:         10000,
				Currency:    "USD",
			},
			expectError: false,
		},
		{
			name: "missing track ID",
			input: CreateLicenseInput{
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{"US"},
				StartDate:   now,
				EndDate:     future,
			},
			expectError: true,
			errorMsg:    "track ID required",
		},
		{
			name: "missing licensee ID",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				RightType:   "sync",
				Territories: []string{"US"},
				StartDate:   now,
				EndDate:     future,
			},
			expectError: true,
			errorMsg:    "licensee ID required",
		},
		{
			name: "invalid dates - end before start",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{"US"},
				StartDate:   future,
				EndDate:     now,
			},
			expectError: true,
			errorMsg:    "end date must be after start date",
		},
		{
			name: "empty territories",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{},
				StartDate:   now,
				EndDate:     future,
			},
			expectError: true,
			errorMsg:    "territories required",
		},
		{
			name: "invalid right type",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "invalid",
				Territories: []string{"US"},
				StartDate:   now,
				EndDate:     future,
			},
			expectError: true,
			errorMsg:    "invalid right type",
		},
		{
			name: "negative fee",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{"US"},
				StartDate:   now,
				EndDate:     future,
				Fee:         -100,
				Currency:    "USD",
			},
			expectError: true,
			errorMsg:    "fee cannot be negative",
		},
		{
			name: "invalid currency",
			input: CreateLicenseInput{
				TrackID:     "track-123",
				LicenseeID:  "licensee-456",
				RightType:   "sync",
				Territories: []string{"US"},
				StartDate:   now,
				EndDate:     future,
				Fee:         10000,
				Currency:    "INVALID",
			},
			expectError: true,
			errorMsg:    "invalid currency code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil) // Will fail - service doesn't exist

			license, err := svc.CreateLicense(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, license)
				assert.NotEmpty(t, license.ID)
				assert.Equal(t, "pending", string(license.Status))
			}
		})
	}
}

func TestLicenseService_GetLicense(t *testing.T) {
	tests := []struct {
		name        string
		licenseID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "existing license",
			licenseID:   "license-123",
			expectError: false,
		},
		{
			name:        "non-existent license",
			licenseID:   "non-existent",
			expectError: true,
			errorMsg:    "license not found",
		},
		{
			name:        "empty license ID",
			licenseID:   "",
			expectError: true,
			errorMsg:    "license ID required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			license, err := svc.GetLicense(ctx, tt.licenseID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, license)
				assert.Equal(t, tt.licenseID, license.ID)
			}
		})
	}
}

func TestLicenseService_ListLicenses(t *testing.T) {
	tests := []struct {
		name        string
		filters     LicenseFilters
		expectCount int
		expectError bool
	}{
		{
			name: "filter by track ID",
			filters: LicenseFilters{
				TrackID: "track-123",
			},
			expectCount: 3, // Example
			expectError: false,
		},
		{
			name: "filter by licensee ID",
			filters: LicenseFilters{
				LicenseeID: "licensee-456",
			},
			expectCount: 2,
			expectError: false,
		},
		{
			name: "filter by status",
			filters: LicenseFilters{
				Status: "active",
			},
			expectCount: 5,
			expectError: false,
		},
		{
			name: "filter by territory",
			filters: LicenseFilters{
				Territory: "US",
			},
			expectCount: 10,
			expectError: false,
		},
		{
			name: "filter by right type",
			filters: LicenseFilters{
				RightType: "sync",
			},
			expectCount: 4,
			expectError: false,
		},
		{
			name: "combined filters",
			filters: LicenseFilters{
				TrackID:   "track-123",
				Status:    "active",
				Territory: "US",
			},
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "no filters - list all",
			filters:     LicenseFilters{},
			expectCount: 20,
			expectError: false,
		},
		{
			name: "invalid status filter",
			filters: LicenseFilters{
				Status: "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			licenses, err := svc.ListLicenses(ctx, tt.filters)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, licenses)
			}
		})
	}
}

func TestLicenseService_UpdateLicense(t *testing.T) {
	now := time.Now()
	newEndDate := now.Add(730 * 24 * time.Hour) // 2 years from now

	tests := []struct {
		name        string
		licenseID   string
		input       UpdateLicenseInput
		expectError bool
		errorMsg    string
	}{
		{
			name:      "update end date",
			licenseID: "license-123",
			input: UpdateLicenseInput{
				EndDate: &newEndDate,
			},
			expectError: false,
		},
		{
			name:      "update auto-renew",
			licenseID: "license-123",
			input: UpdateLicenseInput{
				AutoRenew: boolPtr(true),
			},
			expectError: false,
		},
		{
			name:      "update fee",
			licenseID: "license-123",
			input: UpdateLicenseInput{
				Fee: intPtr(75000),
			},
			expectError: false,
		},
		{
			name:      "update status - pending to active",
			licenseID: "license-pending",
			input: UpdateLicenseInput{
				Status: strPtr("active"),
			},
			expectError: false,
		},
		{
			name:      "invalid status transition - expired to active",
			licenseID: "license-expired",
			input: UpdateLicenseInput{
				Status: strPtr("active"),
			},
			expectError: true,
			errorMsg:    "invalid status transition",
		},
		{
			name:      "invalid status transition - terminated to active",
			licenseID: "license-terminated",
			input: UpdateLicenseInput{
				Status: strPtr("active"),
			},
			expectError: true,
			errorMsg:    "invalid status transition",
		},
		{
			name:        "non-existent license",
			licenseID:   "non-existent",
			input:       UpdateLicenseInput{AutoRenew: boolPtr(true)},
			expectError: true,
			errorMsg:    "license not found",
		},
		{
			name:      "update territories",
			licenseID: "license-123",
			input: UpdateLicenseInput{
				Territories: []string{"US", "CA", "MX"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			license, err := svc.UpdateLicense(ctx, tt.licenseID, tt.input)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, license)
			}
		})
	}
}

func TestLicenseService_RevokeLicense(t *testing.T) {
	tests := []struct {
		name        string
		licenseID   string
		reason      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "revoke active license",
			licenseID:   "license-active",
			reason:      "Contract violation",
			expectError: false,
		},
		{
			name:        "revoke pending license",
			licenseID:   "license-pending",
			reason:      "Deal cancelled",
			expectError: false,
		},
		{
			name:        "revoke already terminated license",
			licenseID:   "license-terminated",
			reason:      "Duplicate",
			expectError: true,
			errorMsg:    "license already terminated",
		},
		{
			name:        "revoke expired license",
			licenseID:   "license-expired",
			reason:      "Cleanup",
			expectError: true,
			errorMsg:    "cannot revoke expired license",
		},
		{
			name:        "non-existent license",
			licenseID:   "non-existent",
			reason:      "Test",
			expectError: true,
			errorMsg:    "license not found",
		},
		{
			name:        "empty reason",
			licenseID:   "license-active",
			reason:      "",
			expectError: true,
			errorMsg:    "reason required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			err := svc.RevokeLicense(ctx, tt.licenseID, tt.reason)

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

func TestLicenseService_GetExpiringLicenses(t *testing.T) {
	tests := []struct {
		name        string
		days        int
		expectCount int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "licenses expiring in 30 days",
			days:        30,
			expectCount: 5,
			expectError: false,
		},
		{
			name:        "licenses expiring in 7 days (for auto-renew)",
			days:        7,
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "licenses expiring in 1 day",
			days:        1,
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "no licenses expiring in 90 days",
			days:        90,
			expectCount: 10, // More licenses in longer timeframe
			expectError: false,
		},
		{
			name:        "negative days",
			days:        -1,
			expectError: true,
			errorMsg:    "days must be positive",
		},
		{
			name:        "zero days",
			days:        0,
			expectError: true,
			errorMsg:    "days must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			licenses, err := svc.GetExpiringLicenses(ctx, tt.days)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, licenses)
			}
		})
	}
}

func TestLicenseService_ProcessAutoRenewals(t *testing.T) {
	tests := []struct {
		name              string
		expectedProcessed int
		expectedFailed    int
		expectError       bool
	}{
		{
			name:              "process renewals successfully",
			expectedProcessed: 3,
			expectedFailed:    0,
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			result, err := svc.ProcessAutoRenewals(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.GreaterOrEqual(t, result.Processed, 0)
			}
		})
	}
}

func TestLicenseService_ActivateLicense(t *testing.T) {
	tests := []struct {
		name        string
		licenseID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "activate pending license",
			licenseID:   "license-pending",
			expectError: false,
		},
		{
			name:        "activate already active license",
			licenseID:   "license-active",
			expectError: true,
			errorMsg:    "license is already active",
		},
		{
			name:        "activate expired license",
			licenseID:   "license-expired",
			expectError: true,
			errorMsg:    "cannot activate expired license",
		},
		{
			name:        "activate terminated license",
			licenseID:   "license-terminated",
			expectError: true,
			errorMsg:    "cannot activate terminated license",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			license, err := svc.ActivateLicense(ctx, tt.licenseID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, license)
				assert.Equal(t, "active", string(license.Status))
			}
		})
	}
}

func TestLicenseService_GetLicensesByHolder(t *testing.T) {
	tests := []struct {
		name        string
		holderID    string
		expectCount int
		expectError bool
	}{
		{
			name:        "holder with multiple licenses",
			holderID:    "holder-123",
			expectCount: 5,
			expectError: false,
		},
		{
			name:        "holder with no licenses",
			holderID:    "holder-no-licenses",
			expectCount: 0,
			expectError: false,
		},
		{
			name:        "empty holder ID",
			holderID:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewLicenseService(nil)

			licenses, err := svc.GetLicensesByHolder(ctx, tt.holderID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, licenses)
			}
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}

// Input types that should be defined in license.go
type CreateLicenseInput struct {
	TrackID     string                 `json:"trackId"`
	LicenseeID  string                 `json:"licenseeId"`
	RightType   string                 `json:"rightType"`
	Territories []string               `json:"territories"`
	StartDate   time.Time              `json:"startDate"`
	EndDate     time.Time              `json:"endDate"`
	AutoRenew   bool                   `json:"autoRenew"`
	Terms       map[string]interface{} `json:"terms,omitempty"`
	Fee         int                    `json:"fee"`
	Currency    string                 `json:"currency"`
}

type UpdateLicenseInput struct {
	EndDate     *time.Time `json:"endDate,omitempty"`
	AutoRenew   *bool      `json:"autoRenew,omitempty"`
	Fee         *int       `json:"fee,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Territories []string   `json:"territories,omitempty"`
}

type LicenseFilters struct {
	TrackID    string `json:"trackId,omitempty"`
	LicenseeID string `json:"licenseeId,omitempty"`
	Status     string `json:"status,omitempty"`
	Territory  string `json:"territory,omitempty"`
	RightType  string `json:"rightType,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type AutoRenewalResult struct {
	Processed int      `json:"processed"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}
