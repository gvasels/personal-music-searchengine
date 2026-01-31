//go:build integration

package testutil

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	s3svc "github.com/aws/aws-sdk-go-v2/service/s3"
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

// ArtistProfileOption allows customizing a test artist profile.
type ArtistProfileOption func(map[string]dynamodbtypes.AttributeValue)

// WithArtistName sets the artist display name.
func WithArtistName(name string) ArtistProfileOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["displayName"] = &dynamodbtypes.AttributeValueMemberS{Value: name}
	}
}

// WithArtistBio sets the artist bio.
func WithArtistBio(bio string) ArtistProfileOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["bio"] = &dynamodbtypes.AttributeValueMemberS{Value: bio}
	}
}

// CreateTestArtistProfile creates an artist profile in DynamoDB and registers it for cleanup.
// PK=USER#{userID}, SK=ARTIST_PROFILE, GSI1PK=ARTIST_PROFILE, GSI1SK=USER#{userID}
// Returns the userID (which is also the profile identifier).
func (tc *TestContext) CreateTestArtistProfile(t *testing.T, userID string, opts ...ArtistProfileOption) string {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "ARTIST_PROFILE"

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":            &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":            &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"GSI1PK":        &dynamodbtypes.AttributeValueMemberS{Value: "ARTIST_PROFILE"},
		"GSI1SK":        &dynamodbtypes.AttributeValueMemberS{Value: "USER#" + userID},
		"Type":          &dynamodbtypes.AttributeValueMemberS{Value: "ARTIST_PROFILE"},
		"userId":        &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"displayName":   &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
		"bio":           &dynamodbtypes.AttributeValueMemberS{Value: "Test bio"},
		"followerCount": &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"trackCount":    &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"albumCount":    &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"totalPlays":    &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"isVerified":    &dynamodbtypes.AttributeValueMemberBOOL{Value: false},
		"createdAt":     &dynamodbtypes.AttributeValueMemberS{Value: now},
		"updatedAt":     &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	for _, opt := range opts {
		opt(item)
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test artist profile: %v", err)
	}

	tc.RegisterCleanup("artist_profile", pk, sk)
	return userID
}

// CreateTestFollow creates a follow relationship in DynamoDB and registers it for cleanup.
// PK=USER#{followerID}, SK=FOLLOWING#{followedID}, GSI1PK=FOLLOWERS#{followedID}, GSI1SK=USER#{followerID}
func (tc *TestContext) CreateTestFollow(t *testing.T, followerID, followedID string) {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + followerID
	sk := "FOLLOWING#" + followedID

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":         &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":         &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"GSI1PK":     &dynamodbtypes.AttributeValueMemberS{Value: "FOLLOWERS#" + followedID},
		"GSI1SK":     &dynamodbtypes.AttributeValueMemberS{Value: "USER#" + followerID},
		"Type":       &dynamodbtypes.AttributeValueMemberS{Value: "FOLLOW"},
		"followerId": &dynamodbtypes.AttributeValueMemberS{Value: followerID},
		"followedId": &dynamodbtypes.AttributeValueMemberS{Value: followedID},
		"createdAt":  &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test follow: %v", err)
	}

	tc.RegisterCleanup("follow", pk, sk)
}

// TagOption allows customizing a test tag.
type TagOption func(map[string]dynamodbtypes.AttributeValue)

// WithTagColor sets the tag color.
func WithTagColor(color string) TagOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["color"] = &dynamodbtypes.AttributeValueMemberS{Value: color}
	}
}

// CreateTestTag creates a tag in DynamoDB and registers it for cleanup.
// PK=USER#{userID}, SK=TAG#{tagName}
// Returns the tag name.
func (tc *TestContext) CreateTestTag(t *testing.T, userID, tagName string, opts ...TagOption) string {
	t.Helper()

	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "TAG#" + tagName

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":         &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":         &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"userId":     &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"name":       &dynamodbtypes.AttributeValueMemberS{Value: tagName},
		"trackCount": &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"CreatedAt":  &dynamodbtypes.AttributeValueMemberS{Value: now},
		"UpdatedAt":  &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	for _, opt := range opts {
		opt(item)
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test tag: %v", err)
	}

	tc.RegisterCleanup("tag", pk, sk)
	return tagName
}

// AlbumOption allows customizing a test album.
type AlbumOption func(map[string]dynamodbtypes.AttributeValue)

// WithAlbumTitle sets the album title.
func WithAlbumTitle(title string) AlbumOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["title"] = &dynamodbtypes.AttributeValueMemberS{Value: title}
	}
}

// WithAlbumArtist sets the album artist.
func WithAlbumArtist(artist string) AlbumOption {
	return func(item map[string]dynamodbtypes.AttributeValue) {
		item["artist"] = &dynamodbtypes.AttributeValueMemberS{Value: artist}
	}
}

// CreateTestAlbum creates an album in DynamoDB and registers it for cleanup.
// PK=USER#{userID}, SK=ALBUM#{albumId}
// Returns the album ID.
func (tc *TestContext) CreateTestAlbum(t *testing.T, userID string, opts ...AlbumOption) string {
	t.Helper()

	albumID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	pk := "USER#" + userID
	sk := "ALBUM#" + albumID

	item := map[string]dynamodbtypes.AttributeValue{
		"PK":            &dynamodbtypes.AttributeValueMemberS{Value: pk},
		"SK":            &dynamodbtypes.AttributeValueMemberS{Value: sk},
		"ID":            &dynamodbtypes.AttributeValueMemberS{Value: albumID},
		"userId":        &dynamodbtypes.AttributeValueMemberS{Value: userID},
		"title":         &dynamodbtypes.AttributeValueMemberS{Value: "Test Album"},
		"artist":        &dynamodbtypes.AttributeValueMemberS{Value: "Test Artist"},
		"trackCount":    &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"totalDuration": &dynamodbtypes.AttributeValueMemberN{Value: "0"},
		"CreatedAt":     &dynamodbtypes.AttributeValueMemberS{Value: now},
		"UpdatedAt":     &dynamodbtypes.AttributeValueMemberS{Value: now},
	}

	for _, opt := range opts {
		opt(item)
	}

	ctx := context.Background()
	_, err := tc.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tc.TableName),
		Item:      item,
	})
	if err != nil {
		t.Fatalf("Failed to create test album: %v", err)
	}

	tc.RegisterCleanup("album", pk, sk)
	return albumID
}

// CreateTestS3Object puts an object in the S3 bucket and registers it for cleanup.
// Returns the key.
func (tc *TestContext) CreateTestS3Object(t *testing.T, key string, content []byte) string {
	t.Helper()

	ctx := context.Background()
	_, err := tc.S3.PutObject(ctx, &s3svc.PutObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		t.Fatalf("Failed to create test S3 object %s: %v", key, err)
	}

	tc.RegisterS3Cleanup(key)
	return key
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
