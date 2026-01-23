package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gvasels/personal-music-searchengine/internal/clients"
	"github.com/gvasels/personal-music-searchengine/internal/handlers"
)

var echoLambda *echoadapter.EchoLambdaV2

func init() {
	// Initialize in init() for Lambda cold start optimization
	if isLambda() {
		e, err := setupEcho()
		if err != nil {
			log.Fatalf("Failed to setup Echo: %v", err)
		}
		echoLambda = echoadapter.NewV2(e)
	}
}

func main() {
	if isLambda() {
		// Run as Lambda
		lambda.Start(echoLambda.ProxyWithContext)
	} else {
		// Run as HTTP server for local development
		e, err := setupEcho()
		if err != nil {
			log.Fatalf("Failed to setup Echo: %v", err)
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "8081"
		}

		log.Printf("Starting Bedrock Gateway server on port %s", port)
		if err := e.Start(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}

func setupEcho() (*echo.Echo, error) {
	ctx := context.Background()

	// Load AWS configuration
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	// Create Bedrock Runtime client
	bedrockClient := bedrockruntime.NewFromConfig(awsCfg)

	// Create clients
	bedrockAPIClient := clients.NewBedrockClient(bedrockClient)
	marengoClient := clients.NewMarengoClient(bedrockClient)

	// Create gateway handler
	gatewayHandler := handlers.NewGatewayHandler(bedrockAPIClient, marengoClient)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API key authentication middleware (optional)
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		e.Use(apiKeyAuth(apiKey))
	}

	// Register gateway routes
	gatewayHandler.RegisterGatewayRoutes(e)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	return e, nil
}

// isLambda returns true if running in AWS Lambda
func isLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" ||
		os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// apiKeyAuth creates middleware for API key authentication
func apiKeyAuth(validKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip health check
			if c.Path() == "/health" {
				return next(c)
			}

			// Check Authorization header
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(401, map[string]interface{}{
					"error": map[string]string{
						"message": "Missing API key",
						"type":    "invalid_request_error",
						"code":    "invalid_api_key",
					},
				})
			}

			// Extract Bearer token
			const bearerPrefix = "Bearer "
			if len(auth) < len(bearerPrefix) || auth[:len(bearerPrefix)] != bearerPrefix {
				return c.JSON(401, map[string]interface{}{
					"error": map[string]string{
						"message": "Invalid API key format",
						"type":    "invalid_request_error",
						"code":    "invalid_api_key",
					},
				})
			}

			key := auth[len(bearerPrefix):]
			if key != validKey {
				return c.JSON(401, map[string]interface{}{
					"error": map[string]string{
						"message": "Invalid API key",
						"type":    "invalid_request_error",
						"code":    "invalid_api_key",
					},
				})
			}

			return next(c)
		}
	}
}
