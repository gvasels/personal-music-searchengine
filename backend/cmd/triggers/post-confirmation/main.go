// Post-Confirmation Lambda Trigger
// Triggered after Cognito user signup confirmation.
// Creates a DynamoDB user profile and adds user to subscriber group.
package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

var (
	userService   service.UserService
	cognitoClient *cognitoidentityprovider.Client
	userPoolID    string
)

func init() {
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		tableName = "music-library"
	}

	// Initialize repository
	repo := repository.NewDynamoDBRepository(dynamoClient, tableName)

	// Initialize user service
	userService = service.NewUserService(repo)

	// Initialize Cognito client
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
	userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	log.Printf("Processing post-confirmation for user: %s", event.UserName)

	// Extract user attributes
	cognitoSub := event.Request.UserAttributes["sub"]
	email := event.Request.UserAttributes["email"]
	name := event.Request.UserAttributes["name"]

	if cognitoSub == "" || email == "" {
		log.Printf("Warning: Missing required attributes for user %s", event.UserName)
		// Return event without error to not block signup
		return event, nil
	}

	// Create user in DynamoDB
	user, err := userService.CreateUserFromCognito(ctx, cognitoSub, email, name)
	if err != nil {
		log.Printf("Error creating user in DynamoDB: %v", err)
		// Return event without error to not block signup
		// User can be created on first API call
		return event, nil
	}

	log.Printf("Created user profile: %s (%s)", user.Email, user.ID)

	// Add user to subscriber group in Cognito
	if userPoolID != "" {
		_, err = cognitoClient.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: &userPoolID,
			Username:   &event.UserName,
			GroupName:  stringPtr("subscriber"),
		})
		if err != nil {
			log.Printf("Warning: Failed to add user to subscriber group: %v", err)
			// Don't fail the signup for this
		} else {
			log.Printf("Added user %s to subscriber group", event.UserName)
		}
	}

	return event, nil
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	lambda.Start(handler)
}
