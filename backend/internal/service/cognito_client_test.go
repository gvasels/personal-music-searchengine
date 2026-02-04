package service

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognitoidentityprovider "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCognitoIdentityProviderAPI mocks the Cognito Identity Provider API
type MockCognitoIdentityProviderAPI struct {
	mock.Mock
}

func (m *MockCognitoIdentityProviderAPI) AdminAddUserToGroup(ctx context.Context, params *cognitoidentityprovider.AdminAddUserToGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminAddUserToGroupOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) AdminRemoveUserFromGroup(ctx context.Context, params *cognitoidentityprovider.AdminRemoveUserFromGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminRemoveUserFromGroupOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) AdminListGroupsForUser(ctx context.Context, params *cognitoidentityprovider.AdminListGroupsForUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminListGroupsForUserOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) AdminDisableUser(ctx context.Context, params *cognitoidentityprovider.AdminDisableUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminDisableUserOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminDisableUserOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) AdminEnableUser(ctx context.Context, params *cognitoidentityprovider.AdminEnableUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminEnableUserOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminEnableUserOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) ListUsers(ctx context.Context, params *cognitoidentityprovider.ListUsersInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUsersOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.ListUsersOutput), args.Error(1)
}

func (m *MockCognitoIdentityProviderAPI) AdminGetUser(ctx context.Context, params *cognitoidentityprovider.AdminGetUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cognitoidentityprovider.AdminGetUserOutput), args.Error(1)
}

func TestCognitoClient_AddUserToGroup(t *testing.T) {
	t.Run("successfully adds user to group", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminAddUserToGroup", ctx, mock.MatchedBy(func(input *cognitoidentityprovider.AdminAddUserToGroupInput) bool {
			return *input.UserPoolId == userPoolID &&
				*input.Username == "user-123" &&
				*input.GroupName == "admin"
		})).Return(&cognitoidentityprovider.AdminAddUserToGroupOutput{}, nil)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.AddUserToGroup(ctx, "user-123", "admin")

		require.NoError(t, err)
		mockAPI.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminAddUserToGroup", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminAddUserToGroupOutput)(nil),
			&types.UserNotFoundException{Message: aws.String("User not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.AddUserToGroup(ctx, "nonexistent", "admin")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("returns error when group not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminAddUserToGroup", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminAddUserToGroupOutput)(nil),
			&types.ResourceNotFoundException{Message: aws.String("Group not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.AddUserToGroup(ctx, "user-123", "nonexistent-group")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("handles AWS service error", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminAddUserToGroup", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminAddUserToGroupOutput)(nil),
			errors.New("service unavailable"),
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.AddUserToGroup(ctx, "user-123", "admin")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "service unavailable")
	})
}

func TestCognitoClient_RemoveUserFromGroup(t *testing.T) {
	t.Run("successfully removes user from group", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminRemoveUserFromGroup", ctx, mock.MatchedBy(func(input *cognitoidentityprovider.AdminRemoveUserFromGroupInput) bool {
			return *input.UserPoolId == userPoolID &&
				*input.Username == "user-123" &&
				*input.GroupName == "admin"
		})).Return(&cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.RemoveUserFromGroup(ctx, "user-123", "admin")

		require.NoError(t, err)
		mockAPI.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminRemoveUserFromGroup", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminRemoveUserFromGroupOutput)(nil),
			&types.UserNotFoundException{Message: aws.String("User not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.RemoveUserFromGroup(ctx, "nonexistent", "admin")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("handles user not in group gracefully", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		// Removing from a group the user isn't in should succeed (idempotent)
		mockAPI.On("AdminRemoveUserFromGroup", ctx, mock.Anything).Return(
			&cognitoidentityprovider.AdminRemoveUserFromGroupOutput{}, nil,
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.RemoveUserFromGroup(ctx, "user-123", "admin")

		require.NoError(t, err)
	})
}

func TestCognitoClient_GetUserGroups(t *testing.T) {
	t.Run("returns user groups", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminListGroupsForUser", ctx, mock.MatchedBy(func(input *cognitoidentityprovider.AdminListGroupsForUserInput) bool {
			return *input.UserPoolId == userPoolID && *input.Username == "user-123"
		})).Return(&cognitoidentityprovider.AdminListGroupsForUserOutput{
			Groups: []types.GroupType{
				{GroupName: aws.String("subscriber")},
				{GroupName: aws.String("artist")},
			},
		}, nil)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		groups, err := client.GetUserGroups(ctx, "user-123")

		require.NoError(t, err)
		assert.Len(t, groups, 2)
		assert.Contains(t, groups, "subscriber")
		assert.Contains(t, groups, "artist")
	})

	t.Run("returns empty slice when user has no groups", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminListGroupsForUser", ctx, mock.Anything).Return(
			&cognitoidentityprovider.AdminListGroupsForUserOutput{
				Groups: []types.GroupType{},
			}, nil,
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		groups, err := client.GetUserGroups(ctx, "user-123")

		require.NoError(t, err)
		assert.Empty(t, groups)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminListGroupsForUser", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminListGroupsForUserOutput)(nil),
			&types.UserNotFoundException{Message: aws.String("User not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		_, err := client.GetUserGroups(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}

func TestCognitoClient_DisableUser(t *testing.T) {
	t.Run("successfully disables user", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminDisableUser", ctx, mock.MatchedBy(func(input *cognitoidentityprovider.AdminDisableUserInput) bool {
			return *input.UserPoolId == userPoolID && *input.Username == "user-123"
		})).Return(&cognitoidentityprovider.AdminDisableUserOutput{}, nil)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.DisableUser(ctx, "user-123")

		require.NoError(t, err)
		mockAPI.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminDisableUser", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminDisableUserOutput)(nil),
			&types.UserNotFoundException{Message: aws.String("User not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.DisableUser(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("disabling already disabled user succeeds", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		// Disabling an already disabled user should be idempotent
		mockAPI.On("AdminDisableUser", ctx, mock.Anything).Return(
			&cognitoidentityprovider.AdminDisableUserOutput{}, nil,
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.DisableUser(ctx, "user-123")

		require.NoError(t, err)
	})
}

func TestCognitoClient_EnableUser(t *testing.T) {
	t.Run("successfully enables user", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminEnableUser", ctx, mock.MatchedBy(func(input *cognitoidentityprovider.AdminEnableUserInput) bool {
			return *input.UserPoolId == userPoolID && *input.Username == "user-123"
		})).Return(&cognitoidentityprovider.AdminEnableUserOutput{}, nil)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.EnableUser(ctx, "user-123")

		require.NoError(t, err)
		mockAPI.AssertExpectations(t)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		mockAPI.On("AdminEnableUser", ctx, mock.Anything).Return(
			(*cognitoidentityprovider.AdminEnableUserOutput)(nil),
			&types.UserNotFoundException{Message: aws.String("User not found")},
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.EnableUser(ctx, "nonexistent")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("enabling already enabled user succeeds", func(t *testing.T) {
		ctx := context.Background()
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		// Enabling an already enabled user should be idempotent
		mockAPI.On("AdminEnableUser", ctx, mock.Anything).Return(
			&cognitoidentityprovider.AdminEnableUserOutput{}, nil,
		)

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		err := client.EnableUser(ctx, "user-123")

		require.NoError(t, err)
	})
}

func TestCognitoClient_Constructor(t *testing.T) {
	t.Run("creates client with user pool ID", func(t *testing.T) {
		mockAPI := new(MockCognitoIdentityProviderAPI)
		userPoolID := "us-east-1_abc123"

		client := NewCognitoClientWithAPI(mockAPI, userPoolID)
		require.NotNil(t, client)
	})
}
