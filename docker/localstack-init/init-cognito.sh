#!/bin/bash

# LocalStack Cognito initialization script
# Creates user pool, app client, groups, and test users for local development

set -e

echo "Initializing LocalStack Cognito resources..."

LOCALSTACK_HOST="${LOCALSTACK_HOST:-localhost}"
AWS_REGION="us-east-1"
USER_POOL_NAME="music-library-local-pool"
CLIENT_NAME="music-library-local-client"
ENDPOINT_URL="http://${LOCALSTACK_HOST}:4566"

# Test user credentials
TEST_PASSWORD="LocalTest123!"

# Wait for LocalStack Cognito to be ready
echo "Waiting for LocalStack Cognito to be ready..."
max_attempts=30
attempt=0
until aws --endpoint-url=${ENDPOINT_URL} cognito-idp list-user-pools --max-results 1 --region ${AWS_REGION} > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "ERROR: LocalStack Cognito not ready after ${max_attempts} attempts"
        exit 1
    fi
    echo "LocalStack Cognito is not ready yet... (attempt ${attempt}/${max_attempts})"
    sleep 2
done
echo "LocalStack Cognito is ready!"

# Check if user pool already exists
EXISTING_POOL=$(aws --endpoint-url=${ENDPOINT_URL} cognito-idp list-user-pools \
    --max-results 10 \
    --region ${AWS_REGION} \
    --query "UserPools[?Name=='${USER_POOL_NAME}'].Id" \
    --output text 2>/dev/null || echo "")

if [ -n "$EXISTING_POOL" ] && [ "$EXISTING_POOL" != "None" ]; then
    echo "User pool ${USER_POOL_NAME} already exists with ID: ${EXISTING_POOL}"
    USER_POOL_ID=$EXISTING_POOL
else
    # Create User Pool
    echo "Creating Cognito User Pool: ${USER_POOL_NAME}"
    USER_POOL_ID=$(aws --endpoint-url=${ENDPOINT_URL} cognito-idp create-user-pool \
        --pool-name ${USER_POOL_NAME} \
        --auto-verified-attributes email \
        --username-attributes email \
        --policies "PasswordPolicy={MinimumLength=8,RequireUppercase=true,RequireLowercase=true,RequireNumbers=true,RequireSymbols=false}" \
        --schema "Name=email,Required=true,Mutable=true" \
        --region ${AWS_REGION} \
        --query 'UserPool.Id' \
        --output text)
    echo "Created User Pool with ID: ${USER_POOL_ID}"
fi

# Check if app client already exists
EXISTING_CLIENT=$(aws --endpoint-url=${ENDPOINT_URL} cognito-idp list-user-pool-clients \
    --user-pool-id ${USER_POOL_ID} \
    --region ${AWS_REGION} \
    --query "UserPoolClients[?ClientName=='${CLIENT_NAME}'].ClientId" \
    --output text 2>/dev/null || echo "")

if [ -n "$EXISTING_CLIENT" ] && [ "$EXISTING_CLIENT" != "None" ]; then
    echo "App client ${CLIENT_NAME} already exists with ID: ${EXISTING_CLIENT}"
    CLIENT_ID=$EXISTING_CLIENT
else
    # Create App Client (no secret for SPA)
    echo "Creating App Client: ${CLIENT_NAME}"
    CLIENT_ID=$(aws --endpoint-url=${ENDPOINT_URL} cognito-idp create-user-pool-client \
        --user-pool-id ${USER_POOL_ID} \
        --client-name ${CLIENT_NAME} \
        --no-generate-secret \
        --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH ALLOW_USER_SRP_AUTH \
        --region ${AWS_REGION} \
        --query 'UserPoolClient.ClientId' \
        --output text)
    echo "Created App Client with ID: ${CLIENT_ID}"
fi

# Create groups (idempotent - will fail silently if exists)
echo "Creating user groups..."
for GROUP in admin artist subscriber; do
    aws --endpoint-url=${ENDPOINT_URL} cognito-idp create-group \
        --user-pool-id ${USER_POOL_ID} \
        --group-name ${GROUP} \
        --description "${GROUP} role group" \
        --region ${AWS_REGION} 2>/dev/null || echo "Group ${GROUP} already exists"
done

# Function to create user and add to group
create_test_user() {
    local EMAIL=$1
    local GROUP=$2

    # Check if user exists
    USER_EXISTS=$(aws --endpoint-url=${ENDPOINT_URL} cognito-idp admin-get-user \
        --user-pool-id ${USER_POOL_ID} \
        --username ${EMAIL} \
        --region ${AWS_REGION} 2>/dev/null && echo "yes" || echo "no")

    if [ "$USER_EXISTS" = "yes" ]; then
        echo "User ${EMAIL} already exists"
    else
        echo "Creating user: ${EMAIL}"

        # Create user
        aws --endpoint-url=${ENDPOINT_URL} cognito-idp admin-create-user \
            --user-pool-id ${USER_POOL_ID} \
            --username ${EMAIL} \
            --user-attributes Name=email,Value=${EMAIL} Name=email_verified,Value=true \
            --message-action SUPPRESS \
            --region ${AWS_REGION}

        # Set password
        aws --endpoint-url=${ENDPOINT_URL} cognito-idp admin-set-user-password \
            --user-pool-id ${USER_POOL_ID} \
            --username ${EMAIL} \
            --password "${TEST_PASSWORD}" \
            --permanent \
            --region ${AWS_REGION}
    fi

    # Add to group (idempotent)
    aws --endpoint-url=${ENDPOINT_URL} cognito-idp admin-add-user-to-group \
        --user-pool-id ${USER_POOL_ID} \
        --username ${EMAIL} \
        --group-name ${GROUP} \
        --region ${AWS_REGION} 2>/dev/null || true

    echo "User ${EMAIL} is in group ${GROUP}"
}

# Create test users
echo ""
echo "Creating test users..."
create_test_user "admin@local.test" "admin"
create_test_user "subscriber@local.test" "subscriber"
create_test_user "artist@local.test" "artist"

echo ""
echo "============================================"
echo "LocalStack Cognito initialization complete!"
echo "============================================"
echo ""
echo "User Pool ID:  ${USER_POOL_ID}"
echo "Client ID:     ${CLIENT_ID}"
echo ""
echo "Test Users (password: ${TEST_PASSWORD}):"
echo "  - admin@local.test      (admin group)"
echo "  - subscriber@local.test (subscriber group)"
echo "  - artist@local.test     (artist group)"
echo ""
echo "Frontend .env.local configuration:"
echo "  VITE_LOCAL_STACK=true"
echo "  VITE_API_URL=http://localhost:8080"
echo "  VITE_COGNITO_USER_POOL_ID=${USER_POOL_ID}"
echo "  VITE_COGNITO_CLIENT_ID=${CLIENT_ID}"
echo "  VITE_COGNITO_ENDPOINT=http://localhost:4566"
echo "  VITE_COGNITO_REGION=us-east-1"
echo ""

# Write config to file for scripts to read
CONFIG_FILE="/tmp/localstack-cognito-config.env"
cat > ${CONFIG_FILE} << EOF
USER_POOL_ID=${USER_POOL_ID}
CLIENT_ID=${CLIENT_ID}
EOF
echo "Configuration saved to ${CONFIG_FILE}"
