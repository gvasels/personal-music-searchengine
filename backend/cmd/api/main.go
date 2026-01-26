package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awslambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gvasels/personal-music-searchengine/internal/handlers"
	"github.com/gvasels/personal-music-searchengine/internal/repository"
	"github.com/gvasels/personal-music-searchengine/internal/search"
	"github.com/gvasels/personal-music-searchengine/internal/service"
)

var echoLambda *echoadapter.EchoLambdaV2

func init() {
	// Initialize in init() for Lambda cold start optimization
	if IsLambda() {
		e, err := setupEcho()
		if err != nil {
			log.Fatalf("Failed to setup Echo: %v", err)
		}
		echoLambda = echoadapter.NewV2(e)
	}
}

func main() {
	if IsLambda() {
		// Run as Lambda
		lambda.Start(echoLambda.ProxyWithContext)
	} else {
		// Run as HTTP server for local development
		e, err := setupEcho()
		if err != nil {
			log.Fatalf("Failed to setup Echo: %v", err)
		}

		cfg, err := LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := e.Start(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}

func setupEcho() (*echo.Echo, error) {
	// Load configuration
	appCfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Load AWS configuration
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(appCfg.AWSRegion))
	if err != nil {
		return nil, err
	}

	// Check for LocalStack endpoint (local development)
	localEndpoint := os.Getenv("AWS_ENDPOINT")

	// Create AWS clients with optional LocalStack endpoint
	var dynamoClient *dynamodb.Client
	var s3Client *s3.Client
	var sfnClient *sfn.Client
	var lambdaClient *awslambda.Client
	var cognitoClient *cognitoidentityprovider.Client

	if localEndpoint != "" {
		// LocalStack configuration
		dynamoClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &localEndpoint
		})
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = &localEndpoint
			o.UsePathStyle = true
		})
		sfnClient = sfn.NewFromConfig(awsCfg, func(o *sfn.Options) {
			o.BaseEndpoint = &localEndpoint
		})
		lambdaClient = awslambda.NewFromConfig(awsCfg, func(o *awslambda.Options) {
			o.BaseEndpoint = &localEndpoint
		})
		cognitoClient = cognitoidentityprovider.NewFromConfig(awsCfg, func(o *cognitoidentityprovider.Options) {
			o.BaseEndpoint = &localEndpoint
		})
	} else {
		dynamoClient = dynamodb.NewFromConfig(awsCfg)
		s3Client = s3.NewFromConfig(awsCfg)
		sfnClient = sfn.NewFromConfig(awsCfg)
		lambdaClient = awslambda.NewFromConfig(awsCfg)
		cognitoClient = cognitoidentityprovider.NewFromConfig(awsCfg)
	}

	// Create repositories
	repo := repository.NewDynamoDBRepository(dynamoClient, appCfg.DynamoDBTableName)
	s3Repo := repository.NewS3Repository(s3Client, s3.NewPresignClient(s3Client), appCfg.MediaBucketName)

	// Create CloudFront signer (optional)
	var cloudfront repository.CloudFrontSigner
	if appCfg.CloudFrontDomain != "" && appCfg.CloudFrontKeyPairID != "" && appCfg.CloudFrontPrivateKey != "" {
		// CloudFront signer would be initialized here
		// For now, we use S3 presigned URLs as fallback
		cloudfront = nil
	}

	// Create services
	services := service.NewServices(
		repo,
		s3Repo,
		cloudfront,
		appCfg.MediaBucketName,
		appCfg.StepFunctionsARN,
	)

	// Set Step Functions client on upload service
	if uploadSvc, ok := services.Upload.(*service.UploadServiceImpl); ok {
		sfnAdapter := service.NewSFNClientAdapter(sfnClient)
		uploadSvc.SetStepFunctionsClient(sfnAdapter)
	}

	// Initialize search service if Nixiesearch function name is configured
	if appCfg.NixiesearchFunctionName != "" {
		searchClient := search.NewClient(lambdaClient, appCfg.NixiesearchFunctionName)
		services.Search = service.NewSearchService(searchClient, repo, s3Repo)
	}

	// Initialize admin service if Cognito User Pool ID is configured
	if appCfg.CognitoUserPoolID != "" {
		cognitoSvc := service.NewCognitoClient(cognitoClient, appCfg.CognitoUserPoolID)
		services.Admin = service.NewAdminService(repo, cognitoSvc)
	}

	// Create handlers
	h := handlers.NewHandlers(services)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.Validator = NewValidator()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Register routes
	h.RegisterRoutes(e)

	// Register admin routes if admin service is configured
	if services.Admin != nil {
		adminHandler := handlers.NewAdminHandler(services.Admin)
		handlers.RegisterAdminRoutes(e, adminHandler)
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	return e, nil
}
