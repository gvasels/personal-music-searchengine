# API Gateway HTTP API with Cognito JWT Authorizer

resource "aws_apigatewayv2_api" "api" {
  name          = "${local.name_prefix}-api"
  protocol_type = "HTTP"
  description   = "Personal Music Search Engine API"

  cors_configuration {
    allow_origins     = ["http://localhost:5173", "http://localhost:3000"]
    allow_methods     = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers     = ["Authorization", "Content-Type", "X-User-ID"]
    expose_headers    = ["X-Request-Id"]
    max_age           = 86400
    allow_credentials = true
  }
}

# Cognito JWT Authorizer
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito"

  jwt_configuration {
    audience = [data.terraform_remote_state.shared.outputs.cognito_client_id]
    issuer   = "https://cognito-idp.${var.aws_region}.amazonaws.com/${local.cognito_user_pool_id}"
  }
}

# Default stage (auto-deployed)
resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      requestTime    = "$context.requestTime"
      httpMethod     = "$context.httpMethod"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      protocol       = "$context.protocol"
      responseLength = "$context.responseLength"
      errorMessage   = "$context.error.message"
    })
  }

  default_route_settings {
    throttling_burst_limit = 100
    throttling_rate_limit  = 50
  }
}

# CloudWatch Log Group for API Gateway
resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "/aws/apigateway/${local.name_prefix}-api"
  retention_in_days = 30
}

# Lambda Integration
resource "aws_apigatewayv2_integration" "api_lambda" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.api.invoke_arn
  payload_format_version = "2.0"
}

# Routes - All routes proxy to Lambda
# User routes
resource "aws_apigatewayv2_route" "get_profile" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/me"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "update_profile" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "PUT /api/v1/me"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Track routes
resource "aws_apigatewayv2_route" "list_tracks" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/tracks"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_track" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/tracks/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "update_track" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "PUT /api/v1/tracks/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "delete_track" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "DELETE /api/v1/tracks/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "add_tags_to_track" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/tracks/{id}/tags"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "remove_tag_from_track" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "DELETE /api/v1/tracks/{id}/tags/{tag}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "upload_cover_art" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "PUT /api/v1/tracks/{id}/cover"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Album routes
resource "aws_apigatewayv2_route" "list_albums" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/albums"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_album" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/albums/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Artist routes
resource "aws_apigatewayv2_route" "list_artists" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/artists"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "list_tracks_by_artist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/artists/{name}/tracks"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "list_albums_by_artist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/artists/{name}/albums"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Playlist routes
resource "aws_apigatewayv2_route" "list_playlists" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/playlists"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "create_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/playlists"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/playlists/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "update_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "PUT /api/v1/playlists/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "delete_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "DELETE /api/v1/playlists/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "add_tracks_to_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/playlists/{id}/tracks"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "remove_tracks_from_playlist" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "DELETE /api/v1/playlists/{id}/tracks"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Tag routes
resource "aws_apigatewayv2_route" "list_tags" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/tags"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "create_tag" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/tags"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_tag" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/tags/{name}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "update_tag" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "PUT /api/v1/tags/{name}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "delete_tag" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "DELETE /api/v1/tags/{name}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_tracks_by_tag" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/tags/{name}/tracks"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Upload routes
resource "aws_apigatewayv2_route" "create_presigned_upload" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/upload/presigned"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "confirm_upload" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/upload/confirm"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "complete_multipart_upload" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/upload/complete-multipart"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "list_uploads" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/uploads"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_upload_status" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/uploads/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "reprocess_upload" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/uploads/{id}/reprocess"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Streaming routes
resource "aws_apigatewayv2_route" "get_stream_url" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/stream/{trackId}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_download_url" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/download/{trackId}"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Search routes
resource "aws_apigatewayv2_route" "simple_search" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "GET /api/v1/search"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "advanced_search" {
  api_id             = aws_apigatewayv2_api.api.id
  route_key          = "POST /api/v1/search"
  target             = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Health check (no auth required)
resource "aws_apigatewayv2_route" "health" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.api_lambda.id}"
}

# Lambda permission for API Gateway
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}
