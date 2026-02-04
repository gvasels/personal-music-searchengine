package models

import (
	"errors"
	"fmt"
	"time"
)

// EntityFollow represents the entity type for follow relationships
const EntityFollow EntityType = "FOLLOW"

// Follow represents a follow relationship between a user and an artist.
type Follow struct {
	FollowerID string    `json:"followerId" dynamodbav:"followerId"`
	FollowedID string    `json:"followedId" dynamodbav:"followedId"`
	CreatedAt  time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

// FollowItem represents a Follow in DynamoDB single-table design
type FollowItem struct {
	DynamoDBItem
	Follow
}

// NewFollow creates a new Follow relationship.
func NewFollow(followerID, followedID string) *Follow {
	return &Follow{
		FollowerID: followerID,
		FollowedID: followedID,
		CreatedAt:  time.Now(),
	}
}

// NewFollowItem creates a DynamoDB item for a follow relationship.
// Primary key pattern: PK=USER#{followerID}, SK=FOLLOWING#{followedID}
// GSI1 pattern: GSI1PK=FOLLOWERS#{followedID}, GSI1SK=USER#{followerID}
func NewFollowItem(follow Follow) FollowItem {
	return FollowItem{
		DynamoDBItem: DynamoDBItem{
			PK:     GetFollowingPK(follow.FollowerID),
			SK:     GetFollowingSK(follow.FollowedID),
			GSI1PK: GetFollowersGSI1PK(follow.FollowedID),
			GSI1SK: fmt.Sprintf("USER#%s", follow.FollowerID),
			Type:   string(EntityFollow),
		},
		Follow: follow,
	}
}

// IsSelfFollow returns true if the follower and followed are the same user.
func (f *Follow) IsSelfFollow() bool {
	return f.FollowerID == f.FollowedID
}

// Validate checks if the follow relationship is valid.
func (f *Follow) Validate() error {
	if f.FollowerID == "" {
		return errors.New("follower ID cannot be empty")
	}
	if f.FollowedID == "" {
		return errors.New("followed ID cannot be empty")
	}
	if f.IsSelfFollow() {
		return errors.New("users cannot follow themselves")
	}
	return nil
}

// GetFollowingPK returns the partition key for querying a user's following list.
func GetFollowingPK(followerID string) string {
	return fmt.Sprintf("USER#%s", followerID)
}

// GetFollowingSK returns the sort key for a follow relationship.
func GetFollowingSK(followedID string) string {
	return fmt.Sprintf("FOLLOWING#%s", followedID)
}

// GetFollowersGSI1PK returns the GSI1 partition key for querying an artist's followers.
func GetFollowersGSI1PK(followedID string) string {
	return fmt.Sprintf("FOLLOWERS#%s", followedID)
}

// FollowResponse represents a follow relationship in API responses.
type FollowResponse struct {
	FollowerID string    `json:"followerId"`
	FollowedID string    `json:"followedId"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ToResponse converts a Follow to a FollowResponse.
func (f *Follow) ToResponse() FollowResponse {
	return FollowResponse{
		FollowerID: f.FollowerID,
		FollowedID: f.FollowedID,
		CreatedAt:  f.CreatedAt,
	}
}

// FollowRequest represents a request to follow an artist.
type FollowRequest struct {
	FollowedID string `json:"followedId" validate:"required,uuid"`
}

// FollowerListResponse represents a list of followers in API responses.
type FollowerListResponse struct {
	Followers  []FollowResponse `json:"followers"`
	TotalCount int              `json:"totalCount"`
	NextKey    string           `json:"nextKey,omitempty"`
}

// FollowingListResponse represents a list of following in API responses.
type FollowingListResponse struct {
	Following  []FollowResponse `json:"following"`
	TotalCount int              `json:"totalCount"`
	NextKey    string           `json:"nextKey,omitempty"`
}
