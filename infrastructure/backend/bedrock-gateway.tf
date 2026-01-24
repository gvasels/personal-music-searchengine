# Bedrock Access Gateway Lambda and API Gateway
# OpenAI-compatible API endpoints for Bedrock models and Marengo video embeddings

# Gateway Lambda Function
resource "aws_lambda_function" "bedrock_gateway" {
  function_name = "${local.name_prefix}-bedrock-gateway"
  role          = aws_iam_role.bedrock_gateway.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # Placeholder - actual code deployed via CI/CD
  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 60

  environment {
    variables = {
      AWS_REGION = var.aws_region
      API_KEY    = aws_secretsmanager_secret.bedrock_gateway_api_key.name
    }
  }

  depends_on = [aws_cloudwatch_log_group.bedrock_gateway_lambda]
}

# CloudWatch Log Group for Gateway Lambda
resource "aws_cloudwatch_log_group" "bedrock_gateway_lambda" {
  name              = "/aws/lambda/${local.name_prefix}-bedrock-gateway"
  retention_in_days = 30
}

# IAM Role for Gateway Lambda
resource "aws_iam_role" "bedrock_gateway" {
  name = "${local.name_prefix}-bedrock-gateway"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# IAM Policy for Bedrock Access
resource "aws_iam_role_policy" "bedrock_gateway_bedrock" {
  name = "bedrock-access"
  role = aws_iam_role.bedrock_gateway.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "bedrock:InvokeModel",
          "bedrock:InvokeModelWithResponseStream"
        ]
        Resource = [
          "arn:aws:bedrock:${var.aws_region}::foundation-model/anthropic.*",
          "arn:aws:bedrock:${var.aws_region}::foundation-model/amazon.*",
          "arn:aws:bedrock:${var.aws_region}::foundation-model/twelvelabs.*"
        ]
      }
    ]
  })
}

# IAM Policy for CloudWatch Logs
resource "aws_iam_role_policy" "bedrock_gateway_logs" {
  name = "cloudwatch-logs"
  role = aws_iam_role.bedrock_gateway.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "${aws_cloudwatch_log_group.bedrock_gateway_lambda.arn}:*"
      }
    ]
  })
}

# IAM Policy for Secrets Manager (API key)
resource "aws_iam_role_policy" "bedrock_gateway_secrets" {
  name = "secrets-access"
  role = aws_iam_role.bedrock_gateway.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = aws_secretsmanager_secret.bedrock_gateway_api_key.arn
      }
    ]
  })
}

# API Key Secret for Gateway Authentication
resource "aws_secretsmanager_secret" "bedrock_gateway_api_key" {
  name                    = "${local.name_prefix}/bedrock-gateway-api-key"
  description             = "API key for Bedrock Gateway authentication"
  recovery_window_in_days = 7
}

# API Gateway HTTP API for Bedrock Gateway
resource "aws_apigatewayv2_api" "bedrock_gateway" {
  name          = "${local.name_prefix}-bedrock-gateway"
  protocol_type = "HTTP"
  description   = "OpenAI-compatible API gateway for Bedrock"

  cors_configuration {
    allow_headers = ["*"]
    allow_methods = ["GET", "POST", "OPTIONS"]
    allow_origins = ["*"]
    max_age       = 3600
  }
}

# API Gateway Integration
resource "aws_apigatewayv2_integration" "bedrock_gateway" {
  api_id           = aws_apigatewayv2_api.bedrock_gateway.id
  integration_type = "AWS_PROXY"

  integration_uri        = aws_lambda_function.bedrock_gateway.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

# API Gateway Routes - OpenAI-compatible endpoints

# POST /v1/chat/completions
resource "aws_apigatewayv2_route" "chat_completions" {
  api_id    = aws_apigatewayv2_api.bedrock_gateway.id
  route_key = "POST /v1/chat/completions"
  target    = "integrations/${aws_apigatewayv2_integration.bedrock_gateway.id}"
}

# POST /v1/embeddings
resource "aws_apigatewayv2_route" "embeddings" {
  api_id    = aws_apigatewayv2_api.bedrock_gateway.id
  route_key = "POST /v1/embeddings"
  target    = "integrations/${aws_apigatewayv2_integration.bedrock_gateway.id}"
}

# POST /v1/embeddings/video (Marengo extension)
resource "aws_apigatewayv2_route" "video_embeddings" {
  api_id    = aws_apigatewayv2_api.bedrock_gateway.id
  route_key = "POST /v1/embeddings/video"
  target    = "integrations/${aws_apigatewayv2_integration.bedrock_gateway.id}"
}

# GET /v1/models
resource "aws_apigatewayv2_route" "models" {
  api_id    = aws_apigatewayv2_api.bedrock_gateway.id
  route_key = "GET /v1/models"
  target    = "integrations/${aws_apigatewayv2_integration.bedrock_gateway.id}"
}

# GET /health
resource "aws_apigatewayv2_route" "bedrock_health" {
  api_id    = aws_apigatewayv2_api.bedrock_gateway.id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.bedrock_gateway.id}"
}

# API Gateway Stage
resource "aws_apigatewayv2_stage" "bedrock_gateway" {
  api_id      = aws_apigatewayv2_api.bedrock_gateway.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.bedrock_gateway_api.arn
    format = jsonencode({
      requestId         = "$context.requestId"
      ip                = "$context.identity.sourceIp"
      requestTime       = "$context.requestTime"
      httpMethod        = "$context.httpMethod"
      routeKey          = "$context.routeKey"
      status            = "$context.status"
      protocol          = "$context.protocol"
      responseLength    = "$context.responseLength"
      integrationError  = "$context.integrationErrorMessage"
      errorMessage      = "$context.error.message"
      latency           = "$context.responseLatency"
    })
  }

  default_route_settings {
    throttling_burst_limit = 100
    throttling_rate_limit  = 50
  }
}

# CloudWatch Log Group for API Gateway access logs
resource "aws_cloudwatch_log_group" "bedrock_gateway_api" {
  name              = "/aws/apigateway/${local.name_prefix}-bedrock-gateway"
  retention_in_days = 30
}

# Lambda Permission for API Gateway
resource "aws_lambda_permission" "bedrock_gateway_api" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bedrock_gateway.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.bedrock_gateway.execution_arn}/*/*"
}

# Outputs
output "bedrock_gateway_api_url" {
  description = "Bedrock Gateway API URL (OpenAI-compatible)"
  value       = aws_apigatewayv2_api.bedrock_gateway.api_endpoint
}

output "bedrock_gateway_lambda_arn" {
  description = "Bedrock Gateway Lambda ARN"
  value       = aws_lambda_function.bedrock_gateway.arn
}

output "bedrock_gateway_api_key_secret" {
  description = "Bedrock Gateway API key secret ARN"
  value       = aws_secretsmanager_secret.bedrock_gateway_api_key.arn
}
