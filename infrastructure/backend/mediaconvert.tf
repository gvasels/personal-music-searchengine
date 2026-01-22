# MediaConvert Infrastructure for HLS Transcoding
# Converts uploaded audio files to adaptive bitrate HLS streams

# MediaConvert Queue (on-demand pricing)
resource "aws_media_convert_queue" "default" {
  name        = "${local.name_prefix}-transcoding"
  pricing_plan = "ON_DEMAND"
  status       = "ACTIVE"
}

# IAM Role for MediaConvert Jobs
resource "aws_iam_role" "mediaconvert" {
  name = "${local.name_prefix}-mediaconvert"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "mediaconvert.amazonaws.com"
        }
      }
    ]
  })
}

# IAM Policy for MediaConvert to access S3
resource "aws_iam_role_policy" "mediaconvert_s3" {
  name = "${local.name_prefix}-mediaconvert-s3"
  role = aws_iam_role.mediaconvert.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:GetObjectVersion"
        ]
        Resource = [
          "${local.media_bucket_arn}/media/*",
          "${local.media_bucket_arn}/uploads/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject"
        ]
        Resource = [
          "${local.media_bucket_arn}/hls/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "*"
      }
    ]
  })
}

# Transcode Start Lambda
resource "aws_lambda_function" "transcode_start" {
  function_name = "${local.name_prefix}-transcode-start"
  role          = aws_iam_role.transcode_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 256
  timeout     = 30

  environment {
    variables = {
      DYNAMODB_TABLE_NAME     = local.dynamodb_table_name
      MEDIA_BUCKET            = local.media_bucket_name
      MEDIACONVERT_ROLE_ARN   = aws_iam_role.mediaconvert.arn
      MEDIACONVERT_QUEUE_ARN  = aws_media_convert_queue.default.arn
      MEDIACONVERT_ENDPOINT   = "https://mediaconvert.${var.aws_region}.amazonaws.com"
    }
  }

  depends_on = [aws_cloudwatch_log_group.transcode_start]
}

resource "aws_cloudwatch_log_group" "transcode_start" {
  name              = "/aws/lambda/${local.name_prefix}-transcode-start"
  retention_in_days = 30
}

# Transcode Complete Lambda (triggered by EventBridge)
resource "aws_lambda_function" "transcode_complete" {
  function_name = "${local.name_prefix}-transcode-complete"
  role          = aws_iam_role.transcode_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 256
  timeout     = 30

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.transcode_complete]
}

resource "aws_cloudwatch_log_group" "transcode_complete" {
  name              = "/aws/lambda/${local.name_prefix}-transcode-complete"
  retention_in_days = 30
}

# IAM Role for Transcode Lambdas
resource "aws_iam_role" "transcode_lambda" {
  name = "${local.name_prefix}-transcode-lambda"

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

# IAM Policy for Transcode Lambdas
resource "aws_iam_role_policy" "transcode_lambda" {
  name = "${local.name_prefix}-transcode-lambda"
  role = aws_iam_role.transcode_lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:UpdateItem"
        ]
        Resource = [
          local.dynamodb_table_arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "mediaconvert:CreateJob",
          "mediaconvert:GetJob",
          "mediaconvert:DescribeEndpoints",
          "mediaconvert:TagResource"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "iam:PassRole"
        ]
        Resource = aws_iam_role.mediaconvert.arn
      }
    ]
  })
}

# Get MediaConvert endpoint for the region
data "aws_region" "current" {}

# Outputs
output "mediaconvert_queue_arn" {
  description = "MediaConvert queue ARN"
  value       = aws_media_convert_queue.default.arn
}

output "transcode_start_lambda_arn" {
  description = "Transcode start Lambda ARN"
  value       = aws_lambda_function.transcode_start.arn
}

output "transcode_complete_lambda_arn" {
  description = "Transcode complete Lambda ARN"
  value       = aws_lambda_function.transcode_complete.arn
}
