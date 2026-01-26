#!/bin/bash
# bootstrap-admin.sh - Promote a user to admin role
#
# Usage: ./bootstrap-admin.sh <email>
#
# This script:
# 1. Looks up the user by email in Cognito
# 2. Adds the user to the 'admin' Cognito group
# 3. Updates the user's role in DynamoDB to 'admin'

set -euo pipefail

# Configuration
AWS_PROFILE="${AWS_PROFILE:-gvasels-muza}"
AWS_REGION="${AWS_REGION:-us-east-1}"
USER_POOL_NAME="music-library-prod-users"
DYNAMODB_TABLE="MusicLibrary"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

usage() {
    echo "Usage: $0 <email>"
    echo ""
    echo "Promotes a user to admin role in both Cognito and DynamoDB."
    echo ""
    echo "Environment variables:"
    echo "  AWS_PROFILE  - AWS profile to use (default: gvasels-muza)"
    echo "  AWS_REGION   - AWS region (default: us-east-1)"
    exit 1
}

# Check arguments
if [ $# -lt 1 ]; then
    usage
fi

EMAIL="$1"

# Validate email format (basic check)
if [[ ! "$EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
    log_error "Invalid email format: $EMAIL"
    exit 1
fi

log_info "Bootstrapping admin user: $EMAIL"
log_info "Using AWS profile: $AWS_PROFILE"
log_info "Region: $AWS_REGION"

# Step 1: Get the User Pool ID
log_info "Looking up User Pool ID..."
USER_POOL_ID=$(aws cognito-idp list-user-pools \
    --max-results 50 \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    --query "UserPools[?Name=='$USER_POOL_NAME'].Id" \
    --output text)

if [ -z "$USER_POOL_ID" ] || [ "$USER_POOL_ID" == "None" ]; then
    log_error "User Pool '$USER_POOL_NAME' not found"
    exit 1
fi

log_info "Found User Pool: $USER_POOL_ID"

# Step 2: Look up user by email
log_info "Looking up user by email..."
USER_INFO=$(aws cognito-idp list-users \
    --user-pool-id "$USER_POOL_ID" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    --filter "email = \"$EMAIL\"" \
    --query "Users[0]" \
    --output json 2>/dev/null || echo "null")

if [ "$USER_INFO" == "null" ] || [ -z "$USER_INFO" ]; then
    log_error "User with email '$EMAIL' not found in Cognito"
    exit 1
fi

# Extract user ID (sub attribute)
USER_SUB=$(echo "$USER_INFO" | jq -r '.Attributes[] | select(.Name == "sub") | .Value')
USERNAME=$(echo "$USER_INFO" | jq -r '.Username')

if [ -z "$USER_SUB" ] || [ "$USER_SUB" == "null" ]; then
    log_error "Could not extract user ID (sub) from user info"
    exit 1
fi

log_info "Found user: $USERNAME (ID: $USER_SUB)"

# Step 3: Add user to admin group in Cognito
log_info "Adding user to admin group in Cognito..."
aws cognito-idp admin-add-user-to-group \
    --user-pool-id "$USER_POOL_ID" \
    --username "$USERNAME" \
    --group-name "admin" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION"

log_info "User added to Cognito 'admin' group"

# Step 4: Update user role in DynamoDB
log_info "Updating user role in DynamoDB..."

# First, check if the user exists in DynamoDB
EXISTING_USER=$(aws dynamodb get-item \
    --table-name "$DYNAMODB_TABLE" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    --key "{\"PK\": {\"S\": \"USER#$USER_SUB\"}, \"SK\": {\"S\": \"PROFILE\"}}" \
    --query "Item" \
    --output json 2>/dev/null || echo "null")

CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if [ "$EXISTING_USER" == "null" ] || [ -z "$EXISTING_USER" ]; then
    log_warn "User not found in DynamoDB, creating new profile..."
    # Create a new user profile with admin role
    aws dynamodb put-item \
        --table-name "$DYNAMODB_TABLE" \
        --profile "$AWS_PROFILE" \
        --region "$AWS_REGION" \
        --item "{
            \"PK\": {\"S\": \"USER#$USER_SUB\"},
            \"SK\": {\"S\": \"PROFILE\"},
            \"EntityType\": {\"S\": \"User\"},
            \"UserID\": {\"S\": \"$USER_SUB\"},
            \"Email\": {\"S\": \"$EMAIL\"},
            \"Role\": {\"S\": \"admin\"},
            \"CreatedAt\": {\"S\": \"$CURRENT_TIME\"},
            \"UpdatedAt\": {\"S\": \"$CURRENT_TIME\"}
        }"
else
    # Update existing user's role
    aws dynamodb update-item \
        --table-name "$DYNAMODB_TABLE" \
        --profile "$AWS_PROFILE" \
        --region "$AWS_REGION" \
        --key "{\"PK\": {\"S\": \"USER#$USER_SUB\"}, \"SK\": {\"S\": \"PROFILE\"}}" \
        --update-expression "SET #role = :role, UpdatedAt = :updated" \
        --expression-attribute-names '{"#role": "Role"}' \
        --expression-attribute-values "{\":role\": {\"S\": \"admin\"}, \":updated\": {\"S\": \"$CURRENT_TIME\"}}"
fi

log_info "User role updated to 'admin' in DynamoDB"

# Verify the changes
log_info "Verifying changes..."

# Check Cognito group membership
GROUPS=$(aws cognito-idp admin-list-groups-for-user \
    --user-pool-id "$USER_POOL_ID" \
    --username "$USERNAME" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    --query "Groups[].GroupName" \
    --output text)

# Check DynamoDB role
DB_ROLE=$(aws dynamodb get-item \
    --table-name "$DYNAMODB_TABLE" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    --key "{\"PK\": {\"S\": \"USER#$USER_SUB\"}, \"SK\": {\"S\": \"PROFILE\"}}" \
    --query "Item.Role.S" \
    --output text)

echo ""
log_info "=== Summary ==="
log_info "User: $EMAIL"
log_info "User ID: $USER_SUB"
log_info "Cognito Groups: $GROUPS"
log_info "DynamoDB Role: $DB_ROLE"
echo ""
log_info "Admin bootstrap complete!"
