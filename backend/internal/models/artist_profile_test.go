package models

import (
	"testing"
)

func TestNewArtistProfile(t *testing.T) {
	userID := "user-123"

	profile := NewArtistProfile(userID)

	if profile.UserID != userID {
		t.Errorf("NewArtistProfile().UserID = %v, want %v", profile.UserID, userID)
	}

	if profile.FollowerCount != 0 {
		t.Errorf("NewArtistProfile().FollowerCount = %v, want 0", profile.FollowerCount)
	}

	if profile.TrackCount != 0 {
		t.Errorf("NewArtistProfile().TrackCount = %v, want 0", profile.TrackCount)
	}

	if profile.IsVerified {
		t.Error("NewArtistProfile().IsVerified should be false by default")
	}

	if profile.SocialLinks == nil {
		t.Error("NewArtistProfile().SocialLinks should be initialized")
	}

	if profile.Genres == nil {
		t.Error("NewArtistProfile().Genres should be initialized")
	}
}

func TestNewArtistProfileItem(t *testing.T) {
	profile := NewArtistProfile("user-123")
	item := NewArtistProfileItem(*profile)

	expectedPK := "USER#user-123"
	expectedSK := "ARTIST_PROFILE"
	expectedGSI1PK := "ARTIST_PROFILE"

	if item.PK != expectedPK {
		t.Errorf("ArtistProfileItem.PK = %v, want %v", item.PK, expectedPK)
	}

	if item.SK != expectedSK {
		t.Errorf("ArtistProfileItem.SK = %v, want %v", item.SK, expectedSK)
	}

	if item.GSI1PK != expectedGSI1PK {
		t.Errorf("ArtistProfileItem.GSI1PK = %v, want %v", item.GSI1PK, expectedGSI1PK)
	}

	if item.Type != string(EntityArtistProfile) {
		t.Errorf("ArtistProfileItem.Type = %v, want %v", item.Type, EntityArtistProfile)
	}
}

func TestArtistProfile_IncrementFollowerCount(t *testing.T) {
	profile := NewArtistProfile("user-123")

	profile.IncrementFollowerCount()
	if profile.FollowerCount != 1 {
		t.Errorf("FollowerCount after increment = %v, want 1", profile.FollowerCount)
	}

	profile.IncrementFollowerCount()
	if profile.FollowerCount != 2 {
		t.Errorf("FollowerCount after second increment = %v, want 2", profile.FollowerCount)
	}
}

func TestArtistProfile_DecrementFollowerCount(t *testing.T) {
	profile := NewArtistProfile("user-123")
	profile.FollowerCount = 5

	profile.DecrementFollowerCount()
	if profile.FollowerCount != 4 {
		t.Errorf("FollowerCount after decrement = %v, want 4", profile.FollowerCount)
	}

	// Should not go below zero
	profile.FollowerCount = 0
	profile.DecrementFollowerCount()
	if profile.FollowerCount != 0 {
		t.Errorf("FollowerCount should not go below 0, got %v", profile.FollowerCount)
	}
}

func TestArtistProfile_ToResponse(t *testing.T) {
	profile := NewArtistProfile("user-123")
	profile.Bio = "Test bio"
	profile.Genres = []string{"rock", "pop"}
	profile.FollowerCount = 100
	profile.IsVerified = true

	response := profile.ToResponse()

	if response.UserID != profile.UserID {
		t.Errorf("Response.UserID = %v, want %v", response.UserID, profile.UserID)
	}

	if response.Bio != profile.Bio {
		t.Errorf("Response.Bio = %v, want %v", response.Bio, profile.Bio)
	}

	if response.FollowerCount != profile.FollowerCount {
		t.Errorf("Response.FollowerCount = %v, want %v", response.FollowerCount, profile.FollowerCount)
	}

	if response.IsVerified != profile.IsVerified {
		t.Errorf("Response.IsVerified = %v, want %v", response.IsVerified, profile.IsVerified)
	}
}

func TestCreateArtistProfileRequest_Validation(t *testing.T) {
	// Valid request
	req := CreateArtistProfileRequest{
		Bio:    "Artist bio",
		Genres: []string{"rock"},
	}

	if req.Bio == "" {
		t.Error("Bio should be set")
	}
}

func TestUpdateArtistProfileRequest(t *testing.T) {
	bio := "Updated bio"
	req := UpdateArtistProfileRequest{
		Bio: &bio,
	}

	if req.Bio == nil || *req.Bio != bio {
		t.Errorf("Bio = %v, want %v", req.Bio, bio)
	}
}
