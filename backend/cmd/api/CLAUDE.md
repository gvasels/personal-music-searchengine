# Backend API Entrypoint - CLAUDE.md

## Overview

Main entrypoint for the API Lambda function. Initializes AWS clients (DynamoDB, S3, Cognito, Step Functions, Lambda), creates the service/handler layers, configures Echo middleware, and runs as either an AWS Lambda handler or a local HTTP server.

## File Descriptions

| File | Purpose |
|------|---------|
| `main.go` | Application entrypoint: Lambda init, local server startup, Echo setup with all routes |
| `config.go` | Environment variable loading into `Config` struct with validation |
| `validator.go` | Echo-compatible request validator using `go-playground/validator` |

## Key Functions

### main.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `init()` | `func init()` | Lambda cold-start optimization: sets up Echo + adapter if running in Lambda |
| `main()` | `func main()` | Entry: runs as Lambda handler or local HTTP server based on environment |
| `setupEcho()` | `func setupEcho() (*echo.Echo, error)` | Creates Echo instance with all AWS clients, repositories, services, handlers, and middleware |

### config.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `LoadConfig()` | `func LoadConfig() (*Config, error)` | Loads config from env vars; requires `DYNAMODB_TABLE_NAME` and `MEDIA_BUCKET` |
| `IsLambda()` | `func IsLambda() bool` | Returns true if `AWS_LAMBDA_FUNCTION_NAME` is set |
| `getEnvOrDefault(key, default)` | `func getEnvOrDefault(string, string) string` | Env var with fallback |

### validator.go

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewValidator()` | `func NewValidator() *CustomValidator` | Creates Echo validator instance |
| `Validate(i)` | `func (cv *CustomValidator) Validate(interface{}) error` | Validates a struct |

## Config Fields

| Field | Env Var | Required | Default |
|-------|---------|----------|---------|
| `AWSRegion` | `AWS_REGION` | No | `us-east-1` |
| `DynamoDBTableName` | `DYNAMODB_TABLE_NAME` | Yes | -- |
| `MediaBucketName` | `MEDIA_BUCKET` | Yes | -- |
| `StepFunctionsARN` | `STEP_FUNCTIONS_ARN` | No | -- |
| `NixiesearchFunctionName` | `NIXIESEARCH_FUNCTION_NAME` | No | -- |
| `CloudFrontDomain` | `CLOUDFRONT_DOMAIN` | No | -- |
| `CloudFrontKeyPairID` | `CLOUDFRONT_KEY_PAIR_ID` | No | -- |
| `CloudFrontPrivateKey` | `CLOUDFRONT_PRIVATE_KEY` | No | -- |
| `CognitoUserPoolID` | `COGNITO_USER_POOL_ID` | No | -- |
| `ServerPort` | `PORT` | No | `8080` |

## Service Wiring (setupEcho)

`setupEcho()` wires the full dependency graph:

1. AWS clients: DynamoDB, S3, SFN, Lambda, Cognito (with optional LocalStack endpoint via `AWS_ENDPOINT`)
2. Repositories: `DynamoDBRepository`, `S3Repository`
3. Services: `NewServices(repo, s3Repo, cloudfront, bucket, sfnArn)` + optional Search, Admin, Hello
4. Handlers: `NewHandlers(services)`, `NewAdminHandler`, `NewHelloHandler`
5. Routes: `RegisterRoutes`, `RegisterAdminRoutes`, `RegisterHelloRoutes`, `/health`
6. Middleware: Logger, Recover, CORS

## Dependencies

### Internal
- `internal/handlers` - HTTP handlers and route registration
- `internal/repository` - DynamoDB and S3 data access
- `internal/service` - Business logic layer
- `internal/search` - Nixiesearch client

### External
- `github.com/labstack/echo/v4` - HTTP framework
- `github.com/aws/aws-lambda-go` - Lambda runtime
- `github.com/awslabs/aws-lambda-go-api-proxy/echo` - Echo-Lambda adapter
- `github.com/aws/aws-sdk-go-v2` - AWS SDK (DynamoDB, S3, SFN, Lambda, Cognito)
- `github.com/go-playground/validator/v10` - Struct validation

## Usage

```bash
# Run locally
DYNAMODB_TABLE_NAME=MusicLibrary MEDIA_BUCKET=music-library-local-media \
  AWS_ENDPOINT=http://localhost:4566 go run ./cmd/api/

# Build for Lambda (ARM64)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap ./cmd/api/
```
