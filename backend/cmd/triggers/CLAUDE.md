# Cognito Triggers - CLAUDE.md

## Overview

Lambda functions triggered by Cognito user pool events. These handlers respond to authentication lifecycle events to synchronize user data between Cognito and DynamoDB.

## Directory Structure

```
triggers/
└── post-confirmation/    # Triggered after user signup confirmation
```

## Trigger Types

| Trigger | Event | Purpose |
|---------|-------|---------|
| `post-confirmation` | PostConfirmation_ConfirmSignUp | Create DynamoDB user profile, assign default role |

## Architecture

```
Cognito User Pool
       │
       ▼ (PostConfirmation trigger)
┌──────────────────────┐
│  post-confirmation   │
│      Lambda          │
└──────────┬───────────┘
           │
     ┌─────┴─────┐
     ▼           ▼
DynamoDB    Cognito Groups
(user)      (subscriber)
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DYNAMODB_TABLE_NAME` | DynamoDB table name | Yes |
| `COGNITO_USER_POOL_ID` | User pool for group assignment | Yes |

## Error Handling

Triggers follow a **fail-safe** pattern:
- Errors are logged but do not block the Cognito flow
- Missing user profiles can be created on first API call
- Group assignment failures are non-blocking

## Build

```bash
cd backend/cmd/triggers/post-confirmation
GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
```

## Deployment

Triggers are deployed via OpenTofu in `infrastructure/shared/cognito.tf` and linked to the Cognito User Pool.
