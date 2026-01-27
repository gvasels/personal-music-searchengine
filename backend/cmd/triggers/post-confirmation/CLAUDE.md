# Post-Confirmation Trigger - CLAUDE.md

## Overview

Cognito post-confirmation Lambda trigger that creates a DynamoDB user profile when a user confirms their signup. Also assigns the user to the `subscriber` Cognito group for default permissions.

## File Description

| File | Purpose |
|------|---------|
| `main.go` | Lambda handler for PostConfirmation_ConfirmSignUp event |

## Flow

```
1. User confirms signup (email verification)
           │
           ▼
2. Cognito invokes PostConfirmation trigger
           │
           ▼
3. Extract user attributes (sub, email, name)
           │
           ▼
4. Create user profile in DynamoDB via UserService
           │
           ▼
5. Add user to "subscriber" Cognito group
           │
           ▼
6. Return event (allows signup to complete)
```

## Handler Signature

```go
func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error)
```

## User Attributes Extracted

| Attribute | Source | Usage |
|-----------|--------|-------|
| `sub` | `event.Request.UserAttributes["sub"]` | User ID (primary key) |
| `email` | `event.Request.UserAttributes["email"]` | User email |
| `name` | `event.Request.UserAttributes["name"]` | Display name (optional) |

## Error Handling

The handler follows a **fail-safe** pattern:
- Missing attributes: Log warning, return success
- DynamoDB errors: Log error, return success (user created on first API call)
- Cognito group errors: Log warning, return success

This ensures signup is never blocked by downstream failures.

## Dependencies

### Internal
- `internal/repository` - DynamoDB repository
- `internal/service` - UserService for profile creation

### External
- `github.com/aws/aws-lambda-go` - Lambda runtime
- `github.com/aws/aws-sdk-go-v2` - AWS SDK (DynamoDB, Cognito)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DYNAMODB_TABLE_NAME` | DynamoDB table name | `music-library` |
| `COGNITO_USER_POOL_ID` | User pool ID for group ops | Required |

## Build

```bash
GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
zip function.zip bootstrap
```

## Testing

This Lambda is tested via integration tests or manual Cognito signup flow. Unit tests should mock:
- `service.UserService.CreateUserFromCognito`
- `cognitoidentityprovider.Client.AdminAddUserToGroup`
