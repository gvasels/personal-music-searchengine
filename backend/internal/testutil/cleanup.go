//go:build integration

package testutil

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CleanupUser removes a user and all their data from DynamoDB.
func (tc *TestContext) CleanupUser(t *testing.T, userID string) {
	t.Helper()

	ctx := context.Background()
	pk := "USER#" + userID

	// Query all items for this user
	items, err := tc.queryUserItems(ctx, pk)
	if err != nil {
		t.Logf("Cleanup warning: failed to query user items for %s: %v", userID, err)
		return
	}

	// Delete each item
	for _, item := range items {
		pkVal, ok := item["PK"].(*dynamodbtypes.AttributeValueMemberS)
		if !ok {
			continue
		}
		skVal, ok := item["SK"].(*dynamodbtypes.AttributeValueMemberS)
		if !ok {
			continue
		}

		err := tc.deleteItem(ctx, pkVal.Value, skVal.Value)
		if err != nil {
			t.Logf("Cleanup warning: failed to delete %s/%s: %v", pkVal.Value, skVal.Value, err)
		}
	}
}

// CleanupTrack removes a specific track from DynamoDB.
func (tc *TestContext) CleanupTrack(t *testing.T, userID, trackID string) {
	t.Helper()

	ctx := context.Background()
	pk := "USER#" + userID
	sk := "TRACK#" + trackID

	err := tc.deleteItem(ctx, pk, sk)
	if err != nil {
		t.Logf("Cleanup warning: failed to delete track %s/%s: %v", userID, trackID, err)
	}
}

// CleanupPlaylist removes a specific playlist from DynamoDB.
func (tc *TestContext) CleanupPlaylist(t *testing.T, userID, playlistID string) {
	t.Helper()

	ctx := context.Background()
	pk := "USER#" + userID
	sk := "PLAYLIST#" + playlistID

	err := tc.deleteItem(ctx, pk, sk)
	if err != nil {
		t.Logf("Cleanup warning: failed to delete playlist %s/%s: %v", userID, playlistID, err)
	}
}

// CleanupAll removes all test data from DynamoDB.
// Use with caution - only for test reset scenarios.
func (tc *TestContext) CleanupAll(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Scan entire table (only safe for local dev)
	if !strings.Contains(tc.Endpoint, "localhost") && !strings.Contains(tc.Endpoint, "4566") {
		t.Log("Cleanup warning: CleanupAll skipped - not running against LocalStack")
		return
	}

	paginator := dynamodb.NewScanPaginator(tc.DynamoDB, &dynamodb.ScanInput{
		TableName:            aws.String(tc.TableName),
		ProjectionExpression: aws.String("PK, SK"),
	})

	deletedCount := 0
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			t.Logf("Cleanup warning: scan failed: %v", err)
			break
		}

		for _, item := range page.Items {
			pkVal, ok := item["PK"].(*dynamodbtypes.AttributeValueMemberS)
			if !ok {
				continue
			}
			skVal, ok := item["SK"].(*dynamodbtypes.AttributeValueMemberS)
			if !ok {
				continue
			}

			// Only delete test data (USER# prefix)
			if !strings.HasPrefix(pkVal.Value, "USER#") {
				continue
			}

			err := tc.deleteItem(ctx, pkVal.Value, skVal.Value)
			if err != nil {
				t.Logf("Cleanup warning: failed to delete %s/%s: %v", pkVal.Value, skVal.Value, err)
			} else {
				deletedCount++
			}
		}
	}

	t.Logf("CleanupAll: deleted %d items", deletedCount)
}

// RegisterS3Cleanup registers an S3 object key for cleanup at test end.
func (tc *TestContext) RegisterS3Cleanup(key string) {
	tc.cleanupItems = append(tc.cleanupItems, cleanupItem{
		itemType: "s3",
		pk:       key,
		sk:       "",
	})
}

// CleanupS3Object removes a specific S3 object.
func (tc *TestContext) CleanupS3Object(t *testing.T, key string) {
	t.Helper()

	ctx := context.Background()
	err := tc.deleteS3Object(ctx, key)
	if err != nil {
		t.Logf("Cleanup warning: failed to delete S3 object %s: %v", key, err)
	}
}

// deleteS3Object deletes an object from S3.
func (tc *TestContext) deleteS3Object(ctx context.Context, key string) error {
	_, err := tc.S3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(tc.BucketName),
		Key:    aws.String(key),
	})
	return err
}

// queryUserItems returns all DynamoDB items for a user.
func (tc *TestContext) queryUserItems(ctx context.Context, pk string) ([]map[string]dynamodbtypes.AttributeValue, error) {
	result, err := tc.DynamoDB.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(tc.TableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]dynamodbtypes.AttributeValue{
			":pk": &dynamodbtypes.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// GetItem retrieves an item from DynamoDB for verification.
func (tc *TestContext) GetItem(t *testing.T, pk, sk string) map[string]dynamodbtypes.AttributeValue {
	t.Helper()

	ctx := context.Background()
	result, err := tc.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tc.TableName),
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		t.Fatalf("Failed to get item %s/%s: %v", pk, sk, err)
	}

	return result.Item
}

// ItemExists checks if an item exists in DynamoDB.
func (tc *TestContext) ItemExists(t *testing.T, pk, sk string) bool {
	t.Helper()

	ctx := context.Background()
	result, err := tc.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tc.TableName),
		Key: map[string]dynamodbtypes.AttributeValue{
			"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk},
			"SK": &dynamodbtypes.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return false
	}

	return len(result.Item) > 0
}
