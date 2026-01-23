package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// RightType represents the type of music right
type RightType string

const (
	RightMechanical  RightType = "mechanical"
	RightPerformance RightType = "performance"
	RightSync        RightType = "sync"
	RightMaster      RightType = "master"
	RightPrint       RightType = "print"
	RightNeighboring RightType = "neighboring"
)

// validRightTypes contains all valid RightType values
var validRightTypes = map[RightType]bool{
	RightMechanical:  true,
	RightPerformance: true,
	RightSync:        true,
	RightMaster:      true,
	RightPrint:       true,
	RightNeighboring: true,
}

// IsValid returns true if the RightType is valid
func (r RightType) IsValid() bool {
	return validRightTypes[r]
}

// String returns the string representation of RightType
func (r RightType) String() string {
	return string(r)
}

// Scope represents the geographic scope of rights or territories
type Scope string

const (
	ScopeGlobal   Scope = "global"
	ScopeNational Scope = "national"
	ScopeRegional Scope = "regional"
	ScopeLocal    Scope = "local"
)

// validScopes contains all valid Scope values
var validScopes = map[Scope]bool{
	ScopeGlobal:   true,
	ScopeNational: true,
	ScopeRegional: true,
	ScopeLocal:    true,
}

// IsValid returns true if the Scope is valid
func (s Scope) IsValid() bool {
	return validScopes[s]
}

// String returns the string representation of Scope
func (s Scope) String() string {
	return string(s)
}

// HolderType represents the type of rights holder
type HolderType string

const (
	HolderLabel       HolderType = "label"
	HolderPublisher   HolderType = "publisher"
	HolderPRO         HolderType = "pro"
	HolderDistributor HolderType = "distributor"
	HolderArtist      HolderType = "artist"
)

// validHolderTypes contains all valid HolderType values
var validHolderTypes = map[HolderType]bool{
	HolderLabel:       true,
	HolderPublisher:   true,
	HolderPRO:         true,
	HolderDistributor: true,
	HolderArtist:      true,
}

// IsValid returns true if the HolderType is valid
func (h HolderType) IsValid() bool {
	return validHolderTypes[h]
}

// String returns the string representation of HolderType
func (h HolderType) String() string {
	return string(h)
}

// LicenseStatus represents the status of a license
type LicenseStatus string

const (
	LicenseActive     LicenseStatus = "active"
	LicenseExpired    LicenseStatus = "expired"
	LicensePending    LicenseStatus = "pending"
	LicenseTerminated LicenseStatus = "terminated"
)

// validLicenseStatuses contains all valid LicenseStatus values
var validLicenseStatuses = map[LicenseStatus]bool{
	LicenseActive:     true,
	LicenseExpired:    true,
	LicensePending:    true,
	LicenseTerminated: true,
}

// IsValid returns true if the LicenseStatus is valid
func (l LicenseStatus) IsValid() bool {
	return validLicenseStatuses[l]
}

// String returns the string representation of LicenseStatus
func (l LicenseStatus) String() string {
	return string(l)
}

// PaymentInfo represents payment details for a rights holder
type PaymentInfo struct {
	Method        string  `dynamodbav:"method" json:"method"`
	Email         *string `dynamodbav:"email,omitempty" json:"email,omitempty"`
	AccountNumber *string `dynamodbav:"accountNumber,omitempty" json:"accountNumber,omitempty"`
	RoutingNumber *string `dynamodbav:"routingNumber,omitempty" json:"routingNumber,omitempty"`
	BankName      *string `dynamodbav:"bankName,omitempty" json:"bankName,omitempty"`
	Currency      string  `dynamodbav:"currency" json:"currency"`
}

// validPaymentMethods contains all valid payment method values
var validPaymentMethods = map[string]bool{
	"paypal":        true,
	"bank_transfer": true,
	"wire":          true,
	"check":         true,
}

// Validate validates the PaymentInfo
func (p *PaymentInfo) Validate() error {
	if p.Method == "" {
		return fmt.Errorf("payment method is required")
	}
	if !validPaymentMethods[p.Method] {
		return fmt.Errorf("invalid payment method: %s", p.Method)
	}
	if p.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if !IsValidCurrencyCode(p.Currency) {
		return fmt.Errorf("invalid currency code: %s", p.Currency)
	}

	// Validate method-specific fields
	switch p.Method {
	case "paypal":
		if p.Email == nil || *p.Email == "" {
			return fmt.Errorf("email is required for paypal payments")
		}
	case "bank_transfer", "wire":
		if p.AccountNumber == nil || *p.AccountNumber == "" {
			return fmt.Errorf("account number is required for bank transfers")
		}
	}

	return nil
}

// TrackRights represents rights information for a track
type TrackRights struct {
	ID           string                 `dynamodbav:"id" json:"id"`
	TrackID      string                 `dynamodbav:"trackId" json:"trackId"`
	HolderID     string                 `dynamodbav:"holderId" json:"holderId"`
	RightType    RightType              `dynamodbav:"rightType" json:"rightType"`
	SharePercent float64                `dynamodbav:"sharePercent" json:"sharePercent"`
	Territories  []string               `dynamodbav:"territories" json:"territories"`
	StartDate    *time.Time             `dynamodbav:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate      *time.Time             `dynamodbav:"endDate,omitempty" json:"endDate,omitempty"`
	Restrictions map[string]interface{} `dynamodbav:"restrictions,omitempty" json:"restrictions,omitempty"`
	CreatedAt    time.Time              `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time              `dynamodbav:"updatedAt" json:"updatedAt"`
}

// Validate validates the TrackRights
func (tr *TrackRights) Validate() error {
	if tr.ID == "" {
		return fmt.Errorf("id is required")
	}
	if tr.TrackID == "" {
		return fmt.Errorf("trackId is required")
	}
	if tr.HolderID == "" {
		return fmt.Errorf("holderId is required")
	}
	if !tr.RightType.IsValid() {
		return fmt.Errorf("invalid rightType: %s", tr.RightType)
	}
	if tr.SharePercent < 0 || tr.SharePercent > 100 {
		return fmt.Errorf("sharePercent must be between 0 and 100")
	}
	if len(tr.Territories) == 0 {
		return fmt.Errorf("territories is required")
	}

	// Validate date range if both dates are provided
	if tr.StartDate != nil && tr.EndDate != nil {
		if tr.EndDate.Before(*tr.StartDate) {
			return fmt.Errorf("endDate must be after startDate")
		}
	}

	return nil
}

// IsActive returns true if the rights are currently active
func (tr *TrackRights) IsActive() bool {
	now := time.Now()

	// Check start date
	if tr.StartDate != nil && now.Before(*tr.StartDate) {
		return false
	}

	// Check end date
	if tr.EndDate != nil && now.After(*tr.EndDate) {
		return false
	}

	return true
}

// AppliesToTerritory returns true if the rights apply to the given territory
func (tr *TrackRights) AppliesToTerritory(territory string) bool {
	if territory == "" {
		return false
	}
	for _, t := range tr.Territories {
		if t == territory {
			return true
		}
	}
	return false
}

// TrackRightsResponse represents the API response for track rights
type TrackRightsResponse struct {
	ID           string                 `json:"id"`
	TrackID      string                 `json:"trackId"`
	HolderID     string                 `json:"holderId"`
	RightType    string                 `json:"rightType"`
	SharePercent float64                `json:"sharePercent"`
	Territories  []string               `json:"territories"`
	StartDate    *time.Time             `json:"startDate,omitempty"`
	EndDate      *time.Time             `json:"endDate,omitempty"`
	Restrictions map[string]interface{} `json:"restrictions,omitempty"`
	IsActive     bool                   `json:"isActive"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// ToResponse converts TrackRights to TrackRightsResponse
func (tr *TrackRights) ToResponse() TrackRightsResponse {
	return TrackRightsResponse{
		ID:           tr.ID,
		TrackID:      tr.TrackID,
		HolderID:     tr.HolderID,
		RightType:    tr.RightType.String(),
		SharePercent: tr.SharePercent,
		Territories:  tr.Territories,
		StartDate:    tr.StartDate,
		EndDate:      tr.EndDate,
		Restrictions: tr.Restrictions,
		IsActive:     tr.IsActive(),
		CreatedAt:    tr.CreatedAt,
		UpdatedAt:    tr.UpdatedAt,
	}
}

// Territory represents a geographic territory for rights management
type Territory struct {
	Code         string                `dynamodbav:"code" json:"code"`
	Name         string                `dynamodbav:"name" json:"name"`
	Scope        Scope                 `dynamodbav:"scope" json:"scope"`
	ParentCode   *string               `dynamodbav:"parentCode,omitempty" json:"parentCode,omitempty"`
	PROID        *string               `dynamodbav:"proId,omitempty" json:"proId,omitempty"`
	RoyaltyRates map[RightType]float64 `dynamodbav:"royaltyRates,omitempty" json:"royaltyRates,omitempty"`
}

// Validate validates the Territory
func (t *Territory) Validate() error {
	if t.Code == "" {
		return fmt.Errorf("code is required")
	}
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !t.Scope.IsValid() {
		return fmt.Errorf("invalid scope: %s", t.Scope)
	}

	// Validate royalty rates if present
	for rightType, rate := range t.RoyaltyRates {
		if !rightType.IsValid() {
			return fmt.Errorf("invalid rightType in royaltyRates: %s", rightType)
		}
		if rate < 0 {
			return fmt.Errorf("royaltyRates cannot be negative")
		}
	}

	return nil
}

// IsGlobal returns true if the territory has global scope
func (t *Territory) IsGlobal() bool {
	return t.Scope == ScopeGlobal
}

// HasParent returns true if the territory has a parent territory
func (t *Territory) HasParent() bool {
	return t.ParentCode != nil && *t.ParentCode != ""
}

// TerritoryResponse represents the API response for a territory
type TerritoryResponse struct {
	Code         string             `json:"code"`
	Name         string             `json:"name"`
	Scope        string             `json:"scope"`
	ParentCode   *string            `json:"parentCode,omitempty"`
	PROID        *string            `json:"proId,omitempty"`
	RoyaltyRates map[string]float64 `json:"royaltyRates,omitempty"`
	IsGlobal     bool               `json:"isGlobal"`
	HasParent    bool               `json:"hasParent"`
}

// ToResponse converts Territory to TerritoryResponse
func (t *Territory) ToResponse() TerritoryResponse {
	// Convert RoyaltyRates to string keys
	var royaltyRates map[string]float64
	if t.RoyaltyRates != nil {
		royaltyRates = make(map[string]float64)
		for k, v := range t.RoyaltyRates {
			royaltyRates[k.String()] = v
		}
	}

	return TerritoryResponse{
		Code:         t.Code,
		Name:         t.Name,
		Scope:        t.Scope.String(),
		ParentCode:   t.ParentCode,
		PROID:        t.PROID,
		RoyaltyRates: royaltyRates,
		IsGlobal:     t.IsGlobal(),
		HasParent:    t.HasParent(),
	}
}

// RightsHolder represents an entity that holds rights to music
type RightsHolder struct {
	ID          string       `dynamodbav:"id" json:"id"`
	Name        string       `dynamodbav:"name" json:"name"`
	Type        HolderType   `dynamodbav:"type" json:"type"`
	Territories []string     `dynamodbav:"territories" json:"territories"`
	IPINumber   *string      `dynamodbav:"ipiNumber,omitempty" json:"ipiNumber,omitempty"`
	ISNI        *string      `dynamodbav:"isni,omitempty" json:"isni,omitempty"`
	ArtistID    *string      `dynamodbav:"artistId,omitempty" json:"artistId,omitempty"`
	PaymentInfo *PaymentInfo `dynamodbav:"paymentInfo,omitempty" json:"-"`
	IsActive    bool         `dynamodbav:"isActive" json:"isActive"`
	CreatedAt   time.Time    `dynamodbav:"createdAt" json:"createdAt"`
}

// Validate validates the RightsHolder
func (rh *RightsHolder) Validate() error {
	if rh.ID == "" {
		return fmt.Errorf("id is required")
	}
	if rh.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !rh.Type.IsValid() {
		return fmt.Errorf("invalid type: %s", rh.Type)
	}
	if len(rh.Territories) == 0 {
		return fmt.Errorf("territories is required")
	}

	// Validate IPI number if present
	if rh.IPINumber != nil && *rh.IPINumber != "" {
		if !IsValidIPINumber(*rh.IPINumber) {
			return fmt.Errorf("invalid ipiNumber format")
		}
	}

	// Validate ISNI if present
	if rh.ISNI != nil && *rh.ISNI != "" {
		if !IsValidISNI(*rh.ISNI) {
			return fmt.Errorf("invalid isni format")
		}
	}

	// Validate payment info if present
	if rh.PaymentInfo != nil {
		if err := rh.PaymentInfo.Validate(); err != nil {
			return fmt.Errorf("invalid payment info: %w", err)
		}
	}

	return nil
}

// IsLabel returns true if the holder is a label
func (rh *RightsHolder) IsLabel() bool {
	return rh.Type == HolderLabel
}

// IsPublisher returns true if the holder is a publisher
func (rh *RightsHolder) IsPublisher() bool {
	return rh.Type == HolderPublisher
}

// IsPRO returns true if the holder is a PRO
func (rh *RightsHolder) IsPRO() bool {
	return rh.Type == HolderPRO
}

// OperatesInTerritory returns true if the holder operates in the given territory
func (rh *RightsHolder) OperatesInTerritory(territory string) bool {
	for _, t := range rh.Territories {
		if t == territory {
			return true
		}
	}
	return false
}

// RightsHolderResponse represents the API response for a rights holder
type RightsHolderResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Territories []string  `json:"territories"`
	IPINumber   *string   `json:"ipiNumber,omitempty"`
	ISNI        *string   `json:"isni,omitempty"`
	ArtistID    *string   `json:"artistId,omitempty"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ToResponse converts RightsHolder to RightsHolderResponse
func (rh *RightsHolder) ToResponse() RightsHolderResponse {
	return RightsHolderResponse{
		ID:          rh.ID,
		Name:        rh.Name,
		Type:        rh.Type.String(),
		Territories: rh.Territories,
		IPINumber:   rh.IPINumber,
		ISNI:        rh.ISNI,
		ArtistID:    rh.ArtistID,
		IsActive:    rh.IsActive,
		CreatedAt:   rh.CreatedAt,
	}
}

// License represents a license for rights usage
type License struct {
	ID          string                 `dynamodbav:"id" json:"id"`
	TrackID     string                 `dynamodbav:"trackId" json:"trackId"`
	LicenseeID  string                 `dynamodbav:"licenseeId" json:"licenseeId"`
	RightType   RightType              `dynamodbav:"rightType" json:"rightType"`
	Territories []string               `dynamodbav:"territories" json:"territories"`
	StartDate   time.Time              `dynamodbav:"startDate" json:"startDate"`
	EndDate     time.Time              `dynamodbav:"endDate" json:"endDate"`
	AutoRenew   bool                   `dynamodbav:"autoRenew" json:"autoRenew"`
	Terms       map[string]interface{} `dynamodbav:"terms" json:"terms"`
	Fee         int                    `dynamodbav:"fee" json:"fee"`
	Currency    string                 `dynamodbav:"currency" json:"currency"`
	Status      LicenseStatus          `dynamodbav:"status" json:"status"`
	CreatedAt   time.Time              `dynamodbav:"createdAt" json:"createdAt"`
}

// Validate validates the License
func (l *License) Validate() error {
	if l.ID == "" {
		return fmt.Errorf("id is required")
	}
	if l.TrackID == "" {
		return fmt.Errorf("trackId is required")
	}
	if l.LicenseeID == "" {
		return fmt.Errorf("licenseeId is required")
	}
	if !l.RightType.IsValid() {
		return fmt.Errorf("invalid rightType: %s", l.RightType)
	}
	if len(l.Territories) == 0 {
		return fmt.Errorf("territories is required")
	}

	// Validate date range
	if l.EndDate.Before(l.StartDate) {
		return fmt.Errorf("endDate must be after startDate")
	}

	// Validate fee
	if l.Fee < 0 {
		return fmt.Errorf("fee cannot be negative")
	}

	// Validate currency
	if l.Currency != "" && !IsValidCurrencyCode(l.Currency) {
		return fmt.Errorf("invalid currency code: %s", l.Currency)
	}

	// Validate status
	if !l.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", l.Status)
	}

	return nil
}

// IsActive returns true if the license is currently active
func (l *License) IsActive() bool {
	if l.Status != LicenseActive {
		return false
	}

	now := time.Now()
	return !now.Before(l.StartDate) && !now.After(l.EndDate)
}

// IsExpired returns true if the license has expired
func (l *License) IsExpired() bool {
	return l.Status == LicenseExpired || time.Now().After(l.EndDate)
}

// CanAutoRenew returns true if the license can auto-renew
func (l *License) CanAutoRenew() bool {
	return l.AutoRenew && l.Status == LicenseActive && !l.IsExpired()
}

// CoversTerritory returns true if the license covers the given territory
func (l *License) CoversTerritory(territory string) bool {
	for _, t := range l.Territories {
		if t == territory {
			return true
		}
	}
	return false
}

// LicenseResponse represents the API response for a license
type LicenseResponse struct {
	ID           string                 `json:"id"`
	TrackID      string                 `json:"trackId"`
	LicenseeID   string                 `json:"licenseeId"`
	RightType    string                 `json:"rightType"`
	Territories  []string               `json:"territories"`
	StartDate    time.Time              `json:"startDate"`
	EndDate      time.Time              `json:"endDate"`
	AutoRenew    bool                   `json:"autoRenew"`
	Terms        map[string]interface{} `json:"terms"`
	Fee          int                    `json:"fee"`
	Currency     string                 `json:"currency"`
	Status       string                 `json:"status"`
	IsActive     bool                   `json:"isActive"`
	IsExpired    bool                   `json:"isExpired"`
	CanAutoRenew bool                   `json:"canAutoRenew"`
	CreatedAt    time.Time              `json:"createdAt"`
}

// ToResponse converts License to LicenseResponse
func (l *License) ToResponse() LicenseResponse {
	return LicenseResponse{
		ID:           l.ID,
		TrackID:      l.TrackID,
		LicenseeID:   l.LicenseeID,
		RightType:    l.RightType.String(),
		Territories:  l.Territories,
		StartDate:    l.StartDate,
		EndDate:      l.EndDate,
		AutoRenew:    l.AutoRenew,
		Terms:        l.Terms,
		Fee:          l.Fee,
		Currency:     l.Currency,
		Status:       l.Status.String(),
		IsActive:     l.IsActive(),
		IsExpired:    l.IsExpired(),
		CanAutoRenew: l.CanAutoRenew(),
		CreatedAt:    l.CreatedAt,
	}
}

// DynamoDB Item Types

// TrackRightsItem represents a TrackRights item in DynamoDB
type TrackRightsItem struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	EntityType string `dynamodbav:"EntityType"`
	TrackRights
}

// NewTrackRightsItem creates a new TrackRightsItem
func NewTrackRightsItem(userID, rightsID string) TrackRightsItem {
	now := time.Now()
	return TrackRightsItem{
		PK:         fmt.Sprintf("USER#%s", userID),
		SK:         fmt.Sprintf("RIGHTS#%s", rightsID),
		EntityType: "rights",
		TrackRights: TrackRights{
			ID:        rightsID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// GetKey returns the DynamoDB key for the TrackRightsItem
func (tri TrackRightsItem) GetKey() map[string]string {
	return map[string]string{
		"PK": tri.PK,
		"SK": tri.SK,
	}
}

// TerritoryItem represents a Territory item in DynamoDB
type TerritoryItem struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	EntityType string `dynamodbav:"EntityType"`
	Territory
}

// NewTerritoryItem creates a new TerritoryItem
func NewTerritoryItem(userID, code string) TerritoryItem {
	return TerritoryItem{
		PK:         fmt.Sprintf("USER#%s", userID),
		SK:         fmt.Sprintf("TERRITORY#%s", code),
		EntityType: "territory",
		Territory: Territory{
			Code: code,
		},
	}
}

// RightsHolderItem represents a RightsHolder item in DynamoDB
type RightsHolderItem struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	EntityType string `dynamodbav:"EntityType"`
	RightsHolder
}

// NewRightsHolderItem creates a new RightsHolderItem
func NewRightsHolderItem(userID, holderID string) RightsHolderItem {
	now := time.Now()
	return RightsHolderItem{
		PK:         fmt.Sprintf("USER#%s", userID),
		SK:         fmt.Sprintf("HOLDER#%s", holderID),
		EntityType: "rightsholder",
		RightsHolder: RightsHolder{
			ID:        holderID,
			CreatedAt: now,
		},
	}
}

// GetKey returns the DynamoDB key for the RightsHolderItem
func (rhi RightsHolderItem) GetKey() map[string]string {
	return map[string]string{
		"PK": rhi.PK,
		"SK": rhi.SK,
	}
}

// LicenseItem represents a License item in DynamoDB
type LicenseItem struct {
	PK         string `dynamodbav:"PK"`
	SK         string `dynamodbav:"SK"`
	EntityType string `dynamodbav:"EntityType"`
	License
}

// NewLicenseItem creates a new LicenseItem
func NewLicenseItem(userID, licenseID string) LicenseItem {
	now := time.Now()
	return LicenseItem{
		PK:         fmt.Sprintf("USER#%s", userID),
		SK:         fmt.Sprintf("LICENSE#%s", licenseID),
		EntityType: "license",
		License: License{
			ID:        licenseID,
			CreatedAt: now,
		},
	}
}

// GetKey returns the DynamoDB key for the LicenseItem
func (li LicenseItem) GetKey() map[string]string {
	return map[string]string{
		"PK": li.PK,
		"SK": li.SK,
	}
}

// Helper Functions

// ValidateSharePercentagesSum validates that share percentages for a right type sum to 100%
func ValidateSharePercentagesSum(rights []TrackRights, rightType RightType) error {
	if len(rights) == 0 {
		return fmt.Errorf("no rights provided")
	}

	var sum float64
	var count int
	for _, r := range rights {
		if r.RightType == rightType {
			sum += r.SharePercent
			count++
		}
	}

	if count == 0 {
		return fmt.Errorf("no rights of type %s found", rightType)
	}

	// Allow for small floating point errors
	if sum < 99.99 || sum > 100.01 {
		return fmt.Errorf("share percentages for %s must sum to 100%%, got %.2f%%", rightType, sum)
	}

	return nil
}

// ValidateAllSharePercentages validates that all right types have shares summing to 100%
func ValidateAllSharePercentages(rights []TrackRights) error {
	if len(rights) == 0 {
		return nil
	}

	// Group rights by type
	rightsByType := make(map[RightType][]TrackRights)
	for _, r := range rights {
		rightsByType[r.RightType] = append(rightsByType[r.RightType], r)
	}

	// Validate each type
	for rightType := range rightsByType {
		if err := ValidateSharePercentagesSum(rights, rightType); err != nil {
			return err
		}
	}

	return nil
}

// validISOCountryCodes contains valid ISO 3166-1 alpha-2 country codes
var validISOCountryCodes = map[string]bool{
	"US": true, "CA": true, "UK": true, "GB": true, "DE": true, "FR": true,
	"IT": true, "ES": true, "NL": true, "BE": true, "AT": true, "CH": true,
	"AU": true, "NZ": true, "JP": true, "KR": true, "CN": true, "HK": true,
	"SG": true, "TW": true, "TH": true, "MY": true, "ID": true, "PH": true,
	"IN": true, "BR": true, "MX": true, "AR": true, "CL": true, "CO": true,
	"PE": true, "ZA": true, "NG": true, "EG": true, "AE": true, "SA": true,
	"IL": true, "TR": true, "RU": true, "PL": true, "CZ": true, "HU": true,
	"RO": true, "GR": true, "PT": true, "SE": true, "NO": true, "DK": true,
	"FI": true, "IE": true, "WW": true,
}

// IsValidISOCountryCode returns true if the code is a valid ISO 3166-1 alpha-2 country code
func IsValidISOCountryCode(code string) bool {
	if code == "" {
		return false
	}
	// Check for regional codes (e.g., US-CA, CA-ON)
	if strings.Contains(code, "-") {
		parts := strings.Split(code, "-")
		if len(parts) == 2 && len(parts[0]) == 2 && len(parts[1]) >= 2 {
			return validISOCountryCodes[strings.ToUpper(parts[0])]
		}
		return false
	}
	return validISOCountryCodes[strings.ToUpper(code)]
}

// validCurrencyCodes contains valid ISO 4217 currency codes
var validCurrencyCodes = map[string]bool{
	"USD": true, "EUR": true, "GBP": true, "JPY": true, "CHF": true, "CAD": true,
	"AUD": true, "NZD": true, "CNY": true, "HKD": true, "SGD": true, "SEK": true,
	"DKK": true, "NOK": true, "KRW": true, "INR": true, "BRL": true, "MXN": true,
	"ZAR": true, "RUB": true, "TRY": true, "PLN": true, "THB": true, "IDR": true,
	"MYR": true, "PHP": true, "CZK": true, "ILS": true, "AED": true, "SAR": true,
}

// IsValidCurrencyCode returns true if the code is a valid ISO 4217 currency code
func IsValidCurrencyCode(code string) bool {
	if len(code) != 3 {
		return false
	}
	return validCurrencyCodes[strings.ToUpper(code)]
}

// ipiRegex matches valid IPI numbers (9-11 digits)
var ipiRegex = regexp.MustCompile(`^\d{9,11}$`)

// IsValidIPINumber returns true if the IPI number is valid
// IPI (Interested Parties Information) numbers are 9-11 digit identifiers
func IsValidIPINumber(ipi string) bool {
	if ipi == "" {
		return false
	}
	return ipiRegex.MatchString(ipi)
}

// IsValidISNI returns true if the ISNI is valid
// ISNI (International Standard Name Identifier) is a 16-character identifier
// Supports formats: 0000000121212121 or 0000-0001-2121-2121
func IsValidISNI(isni string) bool {
	if isni == "" {
		return false
	}

	// Remove dashes for validation
	cleaned := strings.ReplaceAll(isni, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if len(cleaned) != 16 {
		return false
	}

	// First 15 characters must be digits
	for i := 0; i < 15; i++ {
		if cleaned[i] < '0' || cleaned[i] > '9' {
			return false
		}
	}

	// Last character can be digit or X (check digit)
	lastChar := cleaned[15]
	return (lastChar >= '0' && lastChar <= '9') || lastChar == 'X'
}
