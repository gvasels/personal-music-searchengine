# API Lambda Function

resource "aws_lambda_function" "api" {
  function_name = "${local.name_prefix}-api"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  # Placeholder - actual code deployed via CI/CD
  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 256
  timeout     = 30

  environment {
    variables = {
      DYNAMODB_TABLE_NAME           = local.dynamodb_table_name
      MEDIA_BUCKET                  = local.media_bucket_name
      STEP_FUNCTIONS_ARN            = aws_sfn_state_machine.upload_processor.arn
      NIXIESEARCH_FUNCTION_NAME     = aws_lambda_function.nixiesearch.function_name
      CLOUDFRONT_DOMAIN             = aws_cloudfront_distribution.media.domain_name
      CLOUDFRONT_KEY_PAIR_ID        = aws_cloudfront_public_key.signing.id
      CLOUDFRONT_SIGNING_KEY_SECRET = aws_secretsmanager_secret.cloudfront_signing_key.name
      COGNITO_USER_POOL_ID          = local.cognito_user_pool_id
    }
  }

  depends_on = [aws_cloudwatch_log_group.api_lambda]
}

# CloudWatch Log Group for API Lambda
resource "aws_cloudwatch_log_group" "api_lambda" {
  name              = "/aws/lambda/${local.name_prefix}-api"
  retention_in_days = 30
}

# Placeholder archive for initial deployment
data "archive_file" "placeholder" {
  type        = "zip"
  output_path = "${path.module}/placeholder.zip"

  source {
    content  = "placeholder"
    filename = "bootstrap"
  }
}
