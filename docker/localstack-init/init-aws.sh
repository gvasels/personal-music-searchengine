#!/bin/bash

# LocalStack initialization script
# Creates DynamoDB table and S3 bucket for local development

set -e

echo "Initializing LocalStack resources..."

LOCALSTACK_HOST="localhost"
AWS_REGION="us-east-1"
DYNAMODB_TABLE="MusicLibrary"
MEDIA_BUCKET="music-library-local-media"

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
until aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb list-tables --region ${AWS_REGION} > /dev/null 2>&1; do
    echo "LocalStack is not ready yet..."
    sleep 2
done
echo "LocalStack is ready!"

# Create DynamoDB table
echo "Creating DynamoDB table: ${DYNAMODB_TABLE}"
aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 dynamodb create-table \
    --table-name ${DYNAMODB_TABLE} \
    --attribute-definitions \
        AttributeName=PK,AttributeType=S \
        AttributeName=SK,AttributeType=S \
        AttributeName=GSI1PK,AttributeType=S \
        AttributeName=GSI1SK,AttributeType=S \
    --key-schema \
        AttributeName=PK,KeyType=HASH \
        AttributeName=SK,KeyType=RANGE \
    --global-secondary-indexes \
        "[{
            \"IndexName\": \"GSI1\",
            \"KeySchema\": [
                {\"AttributeName\": \"GSI1PK\", \"KeyType\": \"HASH\"},
                {\"AttributeName\": \"GSI1SK\", \"KeyType\": \"RANGE\"}
            ],
            \"Projection\": {\"ProjectionType\": \"ALL\"}
        }]" \
    --billing-mode PAY_PER_REQUEST \
    --region ${AWS_REGION} \
    2>/dev/null || echo "Table ${DYNAMODB_TABLE} already exists"

# Create S3 media bucket
echo "Creating S3 bucket: ${MEDIA_BUCKET}"
aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 s3 mb s3://${MEDIA_BUCKET} \
    --region ${AWS_REGION} \
    2>/dev/null || echo "Bucket ${MEDIA_BUCKET} already exists"

# Configure CORS for the media bucket
echo "Configuring CORS for ${MEDIA_BUCKET}"
aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 s3api put-bucket-cors \
    --bucket ${MEDIA_BUCKET} \
    --cors-configuration '{
        "CORSRules": [
            {
                "AllowedOrigins": ["http://localhost:5173", "http://localhost:3000"],
                "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
                "AllowedHeaders": ["*"],
                "ExposeHeaders": ["ETag", "Content-Length", "Content-Type"],
                "MaxAgeSeconds": 3600
            }
        ]
    }' \
    --region ${AWS_REGION}

# Create bucket folders
echo "Creating bucket folders..."
echo "" | aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 s3 cp - s3://${MEDIA_BUCKET}/uploads/.keep \
    --region ${AWS_REGION}
echo "" | aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 s3 cp - s3://${MEDIA_BUCKET}/media/.keep \
    --region ${AWS_REGION}
echo "" | aws --endpoint-url=http://${LOCALSTACK_HOST}:4566 s3 cp - s3://${MEDIA_BUCKET}/covers/.keep \
    --region ${AWS_REGION}

echo ""
echo "============================================"
echo "LocalStack initialization complete!"
echo "============================================"
echo ""
echo "DynamoDB Table: ${DYNAMODB_TABLE}"
echo "S3 Media Bucket: ${MEDIA_BUCKET}"
echo ""
echo "Connection endpoints:"
echo "  DynamoDB: http://localhost:4566"
echo "  S3:       http://localhost:4566"
echo ""
echo "AWS CLI usage:"
echo "  aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name ${DYNAMODB_TABLE}"
echo "  aws --endpoint-url=http://localhost:4566 s3 ls s3://${MEDIA_BUCKET}/"
echo ""
