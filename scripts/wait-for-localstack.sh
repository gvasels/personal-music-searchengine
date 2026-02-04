#!/bin/bash

# Wait for LocalStack to be healthy and all required services to be running
# Usage: ./wait-for-localstack.sh [timeout_seconds]

set -e

TIMEOUT=${1:-60}
LOCALSTACK_HOST="${LOCALSTACK_HOST:-localhost}"
HEALTH_URL="http://${LOCALSTACK_HOST}:4566/_localstack/health"
REQUIRED_SERVICES="dynamodb s3 cognito-idp"

echo "Waiting for LocalStack to be healthy (timeout: ${TIMEOUT}s)..."
echo "Health URL: ${HEALTH_URL}"
echo "Required services: ${REQUIRED_SERVICES}"
echo ""

start_time=$(date +%s)
last_status=""

while true; do
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))

    if [ $elapsed -ge $TIMEOUT ]; then
        echo ""
        echo "ERROR: Timeout after ${TIMEOUT} seconds waiting for LocalStack"
        echo "Last status: ${last_status}"
        exit 1
    fi

    # Try to fetch health status
    health_response=$(curl -s "${HEALTH_URL}" 2>/dev/null || echo "")

    if [ -z "$health_response" ]; then
        printf "\r[%3ds] LocalStack not responding...                    " $elapsed
        sleep 2
        continue
    fi

    # Check if all required services are running
    all_ready=true
    status_line=""

    for service in $REQUIRED_SERVICES; do
        # Use grep to check service status (works without jq)
        if echo "$health_response" | grep -q "\"${service}\"[[:space:]]*:[[:space:]]*\"running\"" || \
           echo "$health_response" | grep -q "\"${service}\"[[:space:]]*:[[:space:]]*\"available\""; then
            status_line="${status_line}${service}:OK "
        else
            status_line="${status_line}${service}:WAIT "
            all_ready=false
        fi
    done

    last_status="$status_line"
    printf "\r[%3ds] %s" $elapsed "$status_line"

    if [ "$all_ready" = true ]; then
        echo ""
        echo ""
        echo "LocalStack is ready! All services are running."
        exit 0
    fi

    sleep 2
done
