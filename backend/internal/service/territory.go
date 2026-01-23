package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// PRO represents a Performing Rights Organization
type PRO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Description string `json:"description,omitempty"`
}

// RoyaltyRate represents royalty rates for a territory and right type
type RoyaltyRate struct {
	Territory string  `json:"territory"`
	RightType string  `json:"rightType"`
	Rate      float64 `json:"rate"`
	Currency  string  `json:"currency"`
}

// TerritoryService handles territory management operations
type TerritoryService interface {
	GetTerritory(ctx context.Context, code string) (*models.Territory, error)
	ResolveHierarchy(ctx context.Context, code string) ([]*models.Territory, error)
	GetPROForTerritory(ctx context.Context, code string) (*PRO, error)
	GetRoyaltyRates(ctx context.Context, code string, rightType string) (*RoyaltyRate, error)
	ListTerritories(ctx context.Context, scope string) ([]*models.Territory, error)
	IsSubTerritory(ctx context.Context, child, parent string) (bool, error)
}

// territoryService implements TerritoryService
type territoryService struct {
	territories map[string]*models.Territory
	pros        map[string]*PRO
	rates       map[string]map[string]float64 // territory -> rightType -> rate
}

// NewTerritoryService creates a new territory service
func NewTerritoryService(_ interface{}) TerritoryService {
	svc := &territoryService{
		territories: make(map[string]*models.Territory),
		pros:        make(map[string]*PRO),
		rates:       make(map[string]map[string]float64),
	}
	svc.initializeTerritories()
	svc.initializePROs()
	svc.initializeRates()
	return svc
}

// initializeTerritories sets up the territory hierarchy
func (s *territoryService) initializeTerritories() {
	// Global
	s.territories["WW"] = &models.Territory{
		Code:  "WW",
		Name:  "Worldwide",
		Scope: models.ScopeGlobal,
	}

	// National territories (subset)
	national := map[string]string{
		"US": "United States",
		"UK": "United Kingdom",
		"GB": "United Kingdom",
		"DE": "Germany",
		"FR": "France",
		"CA": "Canada",
		"AU": "Australia",
		"JP": "Japan",
		"KR": "South Korea",
		"BR": "Brazil",
		"MX": "Mexico",
		"ES": "Spain",
		"IT": "Italy",
		"NL": "Netherlands",
		"SE": "Sweden",
		"NO": "Norway",
	}

	ww := "WW"
	for code, name := range national {
		s.territories[code] = &models.Territory{
			Code:       code,
			Name:       name,
			Scope:      models.ScopeNational,
			ParentCode: &ww,
		}
	}

	// Regional territories (US states)
	us := "US"
	regions := map[string]string{
		"US-CA": "California",
		"US-NY": "New York",
		"US-TX": "Texas",
		"US-FL": "Florida",
		"US-TN": "Tennessee",
		"US-GA": "Georgia",
	}

	for code, name := range regions {
		s.territories[code] = &models.Territory{
			Code:       code,
			Name:       name,
			Scope:      models.ScopeRegional,
			ParentCode: &us,
		}
	}

	// Local venues (example)
	usca := "US-CA"
	s.territories["LOC:VENUE123"] = &models.Territory{
		Code:       "LOC:VENUE123",
		Name:       "Example Venue LA",
		Scope:      models.ScopeLocal,
		ParentCode: &usca,
	}
}

// initializePROs sets up PRO mappings
func (s *territoryService) initializePROs() {
	s.pros["US"] = &PRO{ID: "ascap", Name: "ASCAP", Country: "US", Description: "American Society of Composers, Authors and Publishers"}
	s.pros["UK"] = &PRO{ID: "prs", Name: "PRS", Country: "UK", Description: "PRS for Music"}
	s.pros["GB"] = &PRO{ID: "prs", Name: "PRS", Country: "GB", Description: "PRS for Music"}
	s.pros["DE"] = &PRO{ID: "gema", Name: "GEMA", Country: "DE", Description: "Gesellschaft für musikalische Aufführungs- und mechanische Vervielfältigungsrechte"}
	s.pros["FR"] = &PRO{ID: "sacem", Name: "SACEM", Country: "FR", Description: "Société des auteurs, compositeurs et éditeurs de musique"}
	s.pros["CA"] = &PRO{ID: "socan", Name: "SOCAN", Country: "CA", Description: "Society of Composers, Authors and Music Publishers of Canada"}
	s.pros["AU"] = &PRO{ID: "apra", Name: "APRA AMCOS", Country: "AU", Description: "Australasian Performing Right Association"}
	s.pros["JP"] = &PRO{ID: "jasrac", Name: "JASRAC", Country: "JP", Description: "Japanese Society for Rights of Authors, Composers and Publishers"}
}

// initializeRates sets up default royalty rates
func (s *territoryService) initializeRates() {
	// US rates
	s.rates["US"] = map[string]float64{
		"mechanical":  0.091, // 9.1 cents per stream
		"performance": 0.12,  // 12% of revenue
		"sync":        0.0,   // Negotiated
		"master":      0.0,   // Negotiated
		"print":       0.10,  // 10%
	}

	// DE rates
	s.rates["DE"] = map[string]float64{
		"mechanical":  0.08,
		"performance": 0.10,
		"sync":        0.0,
		"master":      0.0,
		"print":       0.08,
	}

	// Default rates for other territories
	defaultRates := map[string]float64{
		"mechanical":  0.08,
		"performance": 0.10,
		"sync":        0.0,
		"master":      0.0,
		"print":       0.08,
	}

	for code := range s.territories {
		if _, exists := s.rates[code]; !exists {
			s.rates[code] = defaultRates
		}
	}
}

// GetTerritory retrieves a territory by code
func (s *territoryService) GetTerritory(ctx context.Context, code string) (*models.Territory, error) {
	if code == "" {
		return nil, fmt.Errorf("territory code required")
	}

	territory, exists := s.territories[code]
	if !exists {
		return nil, fmt.Errorf("territory not found: %s", code)
	}

	return territory, nil
}

// ResolveHierarchy returns the territory hierarchy from specific to global
func (s *territoryService) ResolveHierarchy(ctx context.Context, code string) ([]*models.Territory, error) {
	if code == "" {
		return nil, fmt.Errorf("territory code required")
	}

	territory, exists := s.territories[code]
	if !exists {
		return nil, fmt.Errorf("territory not found: %s", code)
	}

	var hierarchy []*models.Territory
	current := territory

	for current != nil {
		hierarchy = append(hierarchy, current)
		if current.ParentCode == nil || *current.ParentCode == "" {
			break
		}
		parent, exists := s.territories[*current.ParentCode]
		if !exists {
			break
		}
		current = parent
	}

	return hierarchy, nil
}

// GetPROForTerritory returns the PRO for a territory
func (s *territoryService) GetPROForTerritory(ctx context.Context, code string) (*PRO, error) {
	if code == "" {
		return nil, fmt.Errorf("territory code required")
	}

	// First check direct mapping
	if pro, exists := s.pros[code]; exists {
		return pro, nil
	}

	// Try to find through hierarchy
	territory, exists := s.territories[code]
	if !exists {
		return nil, fmt.Errorf("territory not found: %s", code)
	}

	// Walk up hierarchy to find PRO
	current := territory
	for current != nil {
		if pro, exists := s.pros[current.Code]; exists {
			return pro, nil
		}
		if current.ParentCode == nil || *current.ParentCode == "" {
			break
		}
		parent, exists := s.territories[*current.ParentCode]
		if !exists {
			break
		}
		current = parent
	}

	return nil, fmt.Errorf("no PRO found for territory: %s", code)
}

// GetRoyaltyRates returns royalty rates for a territory and right type
func (s *territoryService) GetRoyaltyRates(ctx context.Context, code string, rightType string) (*RoyaltyRate, error) {
	if code == "" {
		return nil, fmt.Errorf("territory code required")
	}

	// Validate right type
	validRightTypes := map[string]bool{
		"mechanical": true, "performance": true, "sync": true,
		"master": true, "print": true, "neighboring": true,
	}
	if !validRightTypes[rightType] {
		return nil, fmt.Errorf("invalid right type: %s", rightType)
	}

	// Get territory rates, walk up hierarchy if not found
	territory, exists := s.territories[code]
	if !exists {
		return nil, fmt.Errorf("territory not found: %s", code)
	}

	// Walk up hierarchy to find rates
	current := territory
	for current != nil {
		if rates, exists := s.rates[current.Code]; exists {
			if rate, hasRate := rates[rightType]; hasRate {
				return &RoyaltyRate{
					Territory: code,
					RightType: rightType,
					Rate:      rate,
					Currency:  "USD",
				}, nil
			}
		}
		if current.ParentCode == nil || *current.ParentCode == "" {
			break
		}
		parent, exists := s.territories[*current.ParentCode]
		if !exists {
			break
		}
		current = parent
	}

	// Return default rate
	return &RoyaltyRate{
		Territory: code,
		RightType: rightType,
		Rate:      0.0,
		Currency:  "USD",
	}, nil
}

// ListTerritories returns territories filtered by scope
func (s *territoryService) ListTerritories(ctx context.Context, scope string) ([]*models.Territory, error) {
	validScopes := map[string]models.Scope{
		"global":   models.ScopeGlobal,
		"national": models.ScopeNational,
		"regional": models.ScopeRegional,
		"local":    models.ScopeLocal,
	}

	targetScope, valid := validScopes[strings.ToLower(scope)]
	if !valid {
		return nil, fmt.Errorf("invalid scope: %s", scope)
	}

	var result []*models.Territory
	for _, territory := range s.territories {
		if territory.Scope == targetScope {
			result = append(result, territory)
		}
	}

	return result, nil
}

// IsSubTerritory checks if child is a sub-territory of parent
func (s *territoryService) IsSubTerritory(ctx context.Context, child, parent string) (bool, error) {
	if child == parent {
		return false, nil
	}

	childTerritory, exists := s.territories[child]
	if !exists {
		return false, fmt.Errorf("child territory not found: %s", child)
	}

	parentTerritory, exists := s.territories[parent]
	if !exists {
		return false, fmt.Errorf("parent territory not found: %s", parent)
	}

	// If parent is global, all other territories are sub-territories
	if parentTerritory.Scope == models.ScopeGlobal {
		return true, nil
	}

	// Walk up hierarchy from child
	current := childTerritory
	for current != nil && current.ParentCode != nil {
		if *current.ParentCode == parent {
			return true, nil
		}
		next, exists := s.territories[*current.ParentCode]
		if !exists {
			break
		}
		current = next
	}

	return false, nil
}
