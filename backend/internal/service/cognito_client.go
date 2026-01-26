package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// CognitoClient provides operations for managing users in Cognito User Pools.
// This interface abstracts AWS Cognito admin operations for user management.
type CognitoClient interface {
	// AddUserToGroup adds a user to a Cognito group.
	AddUserToGroup(ctx context.Context, userID string, groupName string) error

	// RemoveUserFromGroup removes a user from a Cognito group.
	RemoveUserFromGroup(ctx context.Context, userID string, groupName string) error

	// GetUserGroups returns the list of groups a user belongs to.
	GetUserGroups(ctx context.Context, userID string) ([]string, error)

	// DisableUser disables a user in Cognito, preventing them from signing in.
	DisableUser(ctx context.Context, userID string) error

	// EnableUser enables a previously disabled user in Cognito.
	EnableUser(ctx context.Context, userID string) error
}

// CognitoIdentityProviderAPI defines the subset of Cognito operations we use.
// This interface allows for mocking in tests.
type CognitoIdentityProviderAPI interface {
	AdminAddUserToGroup(ctx context.Context, params *cognitoidentityprovider.AdminAddUserToGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminRemoveUserFromGroup(ctx context.Context, params *cognitoidentityprovider.AdminRemoveUserFromGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
	AdminListGroupsForUser(ctx context.Context, params *cognitoidentityprovider.AdminListGroupsForUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
	AdminDisableUser(ctx context.Context, params *cognitoidentityprovider.AdminDisableUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error)
	AdminEnableUser(ctx context.Context, params *cognitoidentityprovider.AdminEnableUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error)
}

// cognitoClient implements CognitoClient using AWS SDK v2.
type cognitoClient struct {
	api        CognitoIdentityProviderAPI
	userPoolID string
}

// NewCognitoClient creates a new CognitoClient with the real AWS SDK client.
func NewCognitoClient(client *cognitoidentityprovider.Client, userPoolID string) CognitoClient {
	return &cognitoClient{
		api:        client,
		userPoolID: userPoolID,
	}
}

// NewCognitoClientWithAPI creates a new CognitoClient with a custom API implementation.
// This is primarily used for testing with mocked APIs.
func NewCognitoClientWithAPI(api CognitoIdentityProviderAPI, userPoolID string) CognitoClient {
	return &cognitoClient{
		api:        api,
		userPoolID: userPoolID,
	}
}

// AddUserToGroup adds a user to a Cognito group.
func (c *cognitoClient) AddUserToGroup(ctx context.Context, userID string, groupName string) error {
	input := &cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(userID),
		GroupName:  aws.String(groupName),
	}

	_, err := c.api.AdminAddUserToGroup(ctx, input)
	if err != nil {
		return c.wrapCognitoError(err, "add user to group")
	}
	return nil
}

// RemoveUserFromGroup removes a user from a Cognito group.
func (c *cognitoClient) RemoveUserFromGroup(ctx context.Context, userID string, groupName string) error {
	input := &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(userID),
		GroupName:  aws.String(groupName),
	}

	_, err := c.api.AdminRemoveUserFromGroup(ctx, input)
	if err != nil {
		return c.wrapCognitoError(err, "remove user from group")
	}
	return nil
}

// GetUserGroups returns the list of groups a user belongs to.
func (c *cognitoClient) GetUserGroups(ctx context.Context, userID string) ([]string, error) {
	input := &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(userID),
	}

	output, err := c.api.AdminListGroupsForUser(ctx, input)
	if err != nil {
		return nil, c.wrapCognitoError(err, "get user groups")
	}

	groups := make([]string, 0, len(output.Groups))
	for _, group := range output.Groups {
		if group.GroupName != nil {
			groups = append(groups, *group.GroupName)
		}
	}
	return groups, nil
}

// DisableUser disables a user in Cognito, preventing them from signing in.
func (c *cognitoClient) DisableUser(ctx context.Context, userID string) error {
	input := &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(userID),
	}

	_, err := c.api.AdminDisableUser(ctx, input)
	if err != nil {
		return c.wrapCognitoError(err, "disable user")
	}
	return nil
}

// EnableUser enables a previously disabled user in Cognito.
func (c *cognitoClient) EnableUser(ctx context.Context, userID string) error {
	input := &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(userID),
	}

	_, err := c.api.AdminEnableUser(ctx, input)
	if err != nil {
		return c.wrapCognitoError(err, "enable user")
	}
	return nil
}

// wrapCognitoError converts AWS Cognito errors to appropriate error messages.
func (c *cognitoClient) wrapCognitoError(err error, operation string) error {
	var userNotFound *types.UserNotFoundException
	if errors.As(err, &userNotFound) {
		return fmt.Errorf("failed to %s: user not found", operation)
	}

	var resourceNotFound *types.ResourceNotFoundException
	if errors.As(err, &resourceNotFound) {
		return fmt.Errorf("failed to %s: resource not found", operation)
	}

	var notAuthorized *types.NotAuthorizedException
	if errors.As(err, &notAuthorized) {
		return fmt.Errorf("failed to %s: not authorized", operation)
	}

	return fmt.Errorf("failed to %s: %w", operation, err)
}
