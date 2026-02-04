package models

import (
	"testing"
)

func TestNewFollow(t *testing.T) {
	followerID := "user-123"
	followedID := "artist-456"

	follow := NewFollow(followerID, followedID)

	if follow.FollowerID != followerID {
		t.Errorf("NewFollow().FollowerID = %v, want %v", follow.FollowerID, followerID)
	}

	if follow.FollowedID != followedID {
		t.Errorf("NewFollow().FollowedID = %v, want %v", follow.FollowedID, followedID)
	}

	if follow.CreatedAt.IsZero() {
		t.Error("NewFollow().CreatedAt should be set")
	}
}

func TestNewFollowItem(t *testing.T) {
	follow := NewFollow("user-123", "artist-456")
	item := NewFollowItem(*follow)

	// Primary key: "who does this user follow?"
	expectedPK := "USER#user-123"
	expectedSK := "FOLLOWING#artist-456"

	if item.PK != expectedPK {
		t.Errorf("FollowItem.PK = %v, want %v", item.PK, expectedPK)
	}

	if item.SK != expectedSK {
		t.Errorf("FollowItem.SK = %v, want %v", item.SK, expectedSK)
	}

	// GSI1: "who follows this artist?"
	expectedGSI1PK := "FOLLOWERS#artist-456"
	expectedGSI1SK := "USER#user-123"

	if item.GSI1PK != expectedGSI1PK {
		t.Errorf("FollowItem.GSI1PK = %v, want %v", item.GSI1PK, expectedGSI1PK)
	}

	if item.GSI1SK != expectedGSI1SK {
		t.Errorf("FollowItem.GSI1SK = %v, want %v", item.GSI1SK, expectedGSI1SK)
	}

	if item.Type != string(EntityFollow) {
		t.Errorf("FollowItem.Type = %v, want %v", item.Type, EntityFollow)
	}
}

func TestFollow_IsSelfFollow(t *testing.T) {
	tests := []struct {
		name       string
		followerID string
		followedID string
		expected   bool
	}{
		{"Same user is self-follow", "user-123", "user-123", true},
		{"Different users is not self-follow", "user-123", "artist-456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			follow := NewFollow(tt.followerID, tt.followedID)
			if got := follow.IsSelfFollow(); got != tt.expected {
				t.Errorf("Follow.IsSelfFollow() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFollow_Validate(t *testing.T) {
	tests := []struct {
		name       string
		followerID string
		followedID string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Valid follow",
			followerID: "user-123",
			followedID: "artist-456",
			wantErr:    false,
		},
		{
			name:       "Empty follower ID",
			followerID: "",
			followedID: "artist-456",
			wantErr:    true,
			errMsg:     "follower ID cannot be empty",
		},
		{
			name:       "Empty followed ID",
			followerID: "user-123",
			followedID: "",
			wantErr:    true,
			errMsg:     "followed ID cannot be empty",
		},
		{
			name:       "Self follow",
			followerID: "user-123",
			followedID: "user-123",
			wantErr:    true,
			errMsg:     "users cannot follow themselves",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			follow := NewFollow(tt.followerID, tt.followedID)
			err := follow.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() should return error for %s", tt.name)
				} else if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetFollowingPK(t *testing.T) {
	pk := GetFollowingPK("user-123")
	expected := "USER#user-123"

	if pk != expected {
		t.Errorf("GetFollowingPK() = %v, want %v", pk, expected)
	}
}

func TestGetFollowingSK(t *testing.T) {
	sk := GetFollowingSK("artist-456")
	expected := "FOLLOWING#artist-456"

	if sk != expected {
		t.Errorf("GetFollowingSK() = %v, want %v", sk, expected)
	}
}

func TestGetFollowersGSI1PK(t *testing.T) {
	gsi1pk := GetFollowersGSI1PK("artist-456")
	expected := "FOLLOWERS#artist-456"

	if gsi1pk != expected {
		t.Errorf("GetFollowersGSI1PK() = %v, want %v", gsi1pk, expected)
	}
}

func TestFollow_ToResponse(t *testing.T) {
	follow := NewFollow("user-123", "artist-456")

	response := follow.ToResponse()

	if response.FollowerID != follow.FollowerID {
		t.Errorf("Response.FollowerID = %v, want %v", response.FollowerID, follow.FollowerID)
	}

	if response.FollowedID != follow.FollowedID {
		t.Errorf("Response.FollowedID = %v, want %v", response.FollowedID, follow.FollowedID)
	}

	if response.CreatedAt != follow.CreatedAt {
		t.Errorf("Response.CreatedAt = %v, want %v", response.CreatedAt, follow.CreatedAt)
	}
}
