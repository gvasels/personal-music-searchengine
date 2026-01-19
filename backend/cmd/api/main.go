package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gvasels/personal-music-searchengine/internal/handlers"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	e := setupEcho()
	echoLambda = echoadapter.New(e)
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		lambda.Start(Handler)
	} else {
		// Running locally
		e := setupEcho()
		log.Fatal(e.Start(":8080"))
	}
}

// Handler is the Lambda handler function
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func setupEcho() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Initialize handlers
	h := handlers.NewHandlers()

	// User routes
	api.GET("/me", h.GetCurrentUser)
	api.PUT("/me", h.UpdateCurrentUser)

	// Upload routes
	api.POST("/upload/presigned", h.GetPresignedUploadURL)
	api.POST("/upload/confirm", h.ConfirmUpload)
	api.GET("/uploads", h.ListUploads)

	// Track routes
	api.GET("/tracks", h.ListTracks)
	api.GET("/tracks/:id", h.GetTrack)
	api.PUT("/tracks/:id", h.UpdateTrack)
	api.DELETE("/tracks/:id", h.DeleteTrack)
	api.POST("/tracks/:id/tags", h.AddTagsToTrack)
	api.DELETE("/tracks/:id/tags/:tag", h.RemoveTagFromTrack)

	// Album routes
	api.GET("/albums", h.ListAlbums)
	api.GET("/albums/:id", h.GetAlbum)

	// Artist routes
	api.GET("/artists", h.ListArtists)
	api.GET("/artists/:name", h.GetArtist)

	// Playlist routes
	api.GET("/playlists", h.ListPlaylists)
	api.POST("/playlists", h.CreatePlaylist)
	api.GET("/playlists/:id", h.GetPlaylist)
	api.PUT("/playlists/:id", h.UpdatePlaylist)
	api.DELETE("/playlists/:id", h.DeletePlaylist)
	api.POST("/playlists/:id/tracks", h.AddTracksToPlaylist)
	api.DELETE("/playlists/:id/tracks", h.RemoveTracksFromPlaylist)
	api.PUT("/playlists/:id/tracks/reorder", h.ReorderPlaylistTracks)

	// Tag routes
	api.GET("/tags", h.ListTags)
	api.POST("/tags", h.CreateTag)
	api.PUT("/tags/:name", h.UpdateTag)
	api.DELETE("/tags/:name", h.DeleteTag)

	// Search routes
	api.GET("/search", h.SearchSimple)
	api.POST("/search", h.SearchAdvanced)
	api.GET("/search/suggest", h.GetSearchSuggestions)

	// Streaming routes
	api.GET("/stream/:trackId", h.GetStreamURL)
	api.GET("/download/:trackId", h.GetDownloadURL)
	api.POST("/playback/record", h.RecordPlayback)

	// Queue routes
	api.GET("/queue", h.GetQueue)
	api.PUT("/queue", h.UpdateQueue)
	api.POST("/queue/action", h.QueueAction)

	return e
}
