# Nixiesearch Lambda - Embedded search engine with S3 index storage
# Pure serverless - no VPC, no EFS

# Nixiesearch Lambda Function (Container Image)
resource "aws_lambda_function" "nixiesearch" {
  function_name = "${local.name_prefix}-nixiesearch"
  role          = aws_iam_role.nixiesearch.arn
  package_type  = "Image"
  image_uri     = "${data.terraform_remote_state.global.outputs.ecr_repository_urls.nixiesearch}:latest"

  memory_size = 1024
  timeout     = 30

  # Extended ephemeral storage for search index (2GB)
  ephemeral_storage {
    size = 2048
  }

  environment {
    variables = {
      SEARCH_INDEX_BUCKET = local.search_indexes_bucket_name
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      INDEX_PATH          = "/tmp/nixiesearch"
    }
  }

  depends_on = [aws_cloudwatch_log_group.nixiesearch]
}

resource "aws_cloudwatch_log_group" "nixiesearch" {
  name              = "/aws/lambda/${local.name_prefix}-nixiesearch"
  retention_in_days = 30
}

# IAM Role for Nixiesearch Lambda
resource "aws_iam_role" "nixiesearch" {
  name = "${local.name_prefix}-nixiesearch"

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

# IAM Policy for Nixiesearch Lambda
resource "aws_iam_role_policy" "nixiesearch" {
  name = "${local.name_prefix}-nixiesearch"
  role = aws_iam_role.nixiesearch.id

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
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          local.search_indexes_bucket_arn,
          "${local.search_indexes_bucket_arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:Query"
        ]
        Resource = [
          local.dynamodb_table_arn,
          "${local.dynamodb_table_arn}/index/*"
        ]
      }
    ]
  })
}

# Allow search indexer Lambda to invoke Nixiesearch
resource "aws_lambda_permission" "nixiesearch_from_indexer" {
  statement_id  = "AllowInvokeFromSearchIndexer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.nixiesearch.function_name
  principal     = "lambda.amazonaws.com"
  source_arn    = aws_lambda_function.search_indexer.arn
}

# Allow API Lambda to invoke Nixiesearch
resource "aws_lambda_permission" "nixiesearch_from_api" {
  statement_id  = "AllowInvokeFromAPI"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.nixiesearch.function_name
  principal     = "lambda.amazonaws.com"
  source_arn    = aws_lambda_function.api.arn
}

# IAM Policy for Lambda base role to invoke Nixiesearch
resource "aws_iam_role_policy" "lambda_invoke_nixiesearch" {
  name = "${local.name_prefix}-lambda-invoke-nixiesearch"
  role = split("/", local.lambda_role_arn)[1] # Extract role name from ARN

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "lambda:InvokeFunction"
        Resource = aws_lambda_function.nixiesearch.arn
      }
    ]
  })
}
