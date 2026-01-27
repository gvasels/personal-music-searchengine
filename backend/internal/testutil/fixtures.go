//go:build integration

package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// TestUser represents a test user configuration.
type TestUser struct {
	Email    string
	Password string
	Role     string
	Groups   []string
}

// TestUsers contains pre-defined test users matching init-cognito.sh.
var TestUsers = map[string]TestUser{
	"admin": {
		Email:    "admin@local.test",
		Password: "LocalTest123!",
		Role:     "admin",
		Groups:   []string{"admin"},
	},
	"subscriber": {
		Email:    "subscriber@local.test",
		Password: "LocalTest123!",
		Role:     "subscriber",
		Groups:   []string{"subscriber"},
	},
	"artist": {
		Email:    "artist@local.test",
		Password: "LocalTest123!",
		Role:     "artist",
		Groups:   []string{"artist"},
	},
}

// AuthResult contains authentication tokens.
type AuthResult struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
	ExpiresIn    int32
}

// AuthenticateTestUser authenticates a test user and returns tokens.
func (tc *TestContext) AuthenticateTestUser(t *testing.T, email, password string) *AuthResult {
	t.Helper()

	if tc.UserPoolID == "" || tc.ClientID == "" {
		t.Skip("Cognito not configured. Run init-cognito.sh first.")
	}

	ctx := context.Background()
	result, err := tc.Cognito.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(tc.ClientID),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	})
	if err != nil {
		t.Fatalf("Failed to authenticate user %s: %v", email, err)
	}

	if result.AuthenticationResult == nil {
		t.Fatalf("No authentication result for user %s", email)
	}

	return &AuthResult{
		AccessToken:  aws.ToString(result.AuthenticationResult.AccessToken),
		IDToken:      aws.ToString(result.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(result.AuthenticationResult.RefreshToken),
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
	}
}

// GetTestUserToken returns an access token for a predefined test user.
func (tc *TestContext) GetTestUserToken(t *testing.T, role string) string {
	t.Helper()

	user, ok := TestUsers[role]
	if !ok {
		t.Fatalf("Unknown test user role: %s", role)
	}

	auth := tc.AuthenticateTestUser(t, user.Email, user.Password)
	return auth.AccessToken
}

// TrackOption allows customizing a test track.
type TrackOption func(map[string]dynamodbtypes.AttributeValue)

// WithTrackTitle sets the track title.
func WithTrackTitle(title string) TrackOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["Title"] = &dynamodbtypes.AttributeValueMemberS{Value: title}
	}
}

// WithTrackArtist sets the track artist.
func WithTrackArtist(artist string) TrackOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["Artist"] = &dynamodbtypes.AttributeValueMemberS{Value: artist}
	}
}

// WithTrackAlbum sets the track album.
func WithTrackAlbum(album string) TrackOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["Album"] = &dynamodbtypes.AttributeValueMemberS{Value: album}
	}
}

// WithTrackVisibility sets the track visibility.
func WithTrackVisibility(visibility string) TrackOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["Visibility"] = &dynamodbtypes.AttributeValueMemberS{Value: visibility}
	}
}

// CreateTestTrack creates a track in DynamoDB and registers it for cleanup.
// Returns the track ID.
func (tc *TestContext) CreateTestTrack(t *testing.T, userID string, opts ...TrackOption) string {
	t.Helper()

	trackID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "TRACK#" + trackID

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":         &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":         &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"ID":         &dynamodbtypes.AttributeValueMemberS{Value: trackID},
		"UserID":     &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"Title":      &dynamodbtypes.AttributeValueMemberS{Value: "Test Track"},
		"Artist":     &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
		"Album":      &dynamodbtypes.AttributeValueMemberS{Value: "Test Album"},
		"Duration":   &dynamodbtypes.AttributeValueMemberN{Value: "180"},
		"Format":     &dynamodbtypes.AttributeValueMemberS{Value: "mp3"},
		"Visibility": &dynamodbtypes.AttributeValueMemberS{Value: "private"},
		"CreatedAt":  &dynamodbtypes.AttributeValueMemberS{Value: now},
		"UpdatedAt":  &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	// Apply options
	for _, opt := range opts {
		opt(item)
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test track: %v", err)
	}

	// Register for cleanup
	tc.RegisterCleanup("track", pk, sk)

	return trackID
}

// CreateTestUser creates a user profile in DynamoDB and registers it for cleanup.
// Returns the user ID.
func (tc *TestContext) CreateTestUser(t *testing.T, email, role string) string {
	t.Helper()

	userID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "PROFILE"

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":          &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":          &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"ID":          &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"Email":       &dynamodbtypes.AttributeValueMemberS{Value: email},
		"DisplayName": &dynamodbtypes.AttributeValueMemberS{Value: "Test User"},
		"Role":        &dynamodbtypes.AttributeValueMemberS{Value: role},
		"CreatedAt":   &dynamodbtypes.AttributeValueMemberS{Value: now},
		"UpdatedAt":   &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Register for cleanup
	tc.RegisterCleanup("user", pk, sk)

	return userID
}

// CreateTestPlaylist creates a playlist in DynamoDB and registers it for cleanup.
// Returns the playlist ID.
func (tc *TestContext) CreateTestPlaylist(t *testing.T, userID, name string) string {
	t.Helper()

	playlistID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "PLAYLIST#" + playlistID

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":          &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":          &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"ID":          &dynamodbtypes.AttributeValueMemberS{Value: playlistID},
		"UserID":      &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"Name":        &dynamodbtypes.AttributeValueMemberS{Value: name},
		"Description": &dynamodbtypes.AttributeValueMemberS{Value: "Test playlist"},
		"Visibility":  &dynamodbtypes.AttributeValueMemberS{Value: "private"},
		"TrackCount":  &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"CreatedAt":   &dynamodbtypes.AttributeValueMemberS{Value: now},
		"UpdatedAt":   &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test playlist: %v", err)
	}

	// Register for cleanup
	tc.RegisterCleanup("playlist", pk, sk)

	return playlistID
}
