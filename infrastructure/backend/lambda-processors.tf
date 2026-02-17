# Step Functions Processor Lambda Functions

data "aws_caller_identity" "current" {}

# Metadata Extractor Lambda
resource "aws_lambda_function" "metadata_extractor" {
  function_name = "${local.name_prefix}-metadata-extractor"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 60

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.metadata_extractor]
}

resource "aws_cloudwatch_log_group" "metadata_extractor" {
  name              = "/aws/lambda/${local.name_prefix}-metadata-extractor"
  retention_in_days = 30
}

# Cover Art Processor Lambda
resource "aws_lambda_function" "cover_art_processor" {
  function_name = "${local.name_prefix}-cover-art-processor"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 60

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.cover_art_processor]
}

resource "aws_cloudwatch_log_group" "cover_art_processor" {
  name              = "/aws/lambda/${local.name_prefix}-cover-art-processor"
  retention_in_days = 30
}

# FFmpeg Lambda layer for audio processing (ARM64)
# Placeholder - actual layer code deployed via CI/CD
resource "aws_lambda_layer_version" "ffmpeg" {
  layer_name               = "${local.name_prefix}-ffmpeg"
  filename                 = data.archive_file.placeholder.output_path
  source_code_hash         = data.archive_file.placeholder.output_base64sha256
  compatible_runtimes      = ["provided.al2023"]
  compatible_architectures = ["arm64"]
  description              = "FFmpeg static binaries for ARM64"
}

# Audio Analyzer Lambda (GenAI genre/mood analysis)
resource "aws_lambda_function" "audio_analyzer" {
  function_name = "${local.name_prefix}-audio-analyzer"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["x86_64"]

  filename         = "${path.module}/../../backend/cmd/processor/audio-analyzer/bootstrap.zip"
  source_code_hash = filebase64sha256("${path.module}/../../backend/cmd/processor/audio-analyzer/bootstrap.zip")

  memory_size = 256
  timeout     = 60

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
      BEDROCK_MODEL_ID    = "anthropic.claude-3-haiku-20240307-v1:0"
    }
  }

  depends_on = [aws_cloudwatch_log_group.audio_analyzer]
}

resource "aws_cloudwatch_log_group" "audio_analyzer" {
  name              = "/aws/lambda/${local.name_prefix}-audio-analyzer"
  retention_in_days = 30
}

# Embedding Generator Lambda (Marengo audio embeddings)
resource "aws_lambda_function" "embedding_generator" {
  function_name = "${local.name_prefix}-embedding-generator"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["x86_64"]

  filename         = "${path.module}/../../backend/cmd/processor/embedding-generator/bootstrap.zip"
  source_code_hash = filebase64sha256("${path.module}/../../backend/cmd/processor/embedding-generator/bootstrap.zip")

  memory_size = 256
  timeout     = 120

  environment {
    variables = {
      MEDIA_BUCKET      = local.media_bucket_name
      VECTOR_BUCKET     = local.vector_bucket_name
      VECTOR_INDEX_NAME = "media-embeddings"
      AWS_ACCOUNT_ID    = data.aws_caller_identity.current.account_id
    }
  }

  depends_on = [aws_cloudwatch_log_group.embedding_generator]
}

resource "aws_cloudwatch_log_group" "embedding_generator" {
  name              = "/aws/lambda/${local.name_prefix}-embedding-generator"
  retention_in_days = 30
}

# Track Updater Lambda (updates DynamoDB with analysis results)
resource "aws_lambda_function" "track_updater" {
  function_name = "${local.name_prefix}-track-updater"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["x86_64"]

  filename         = "${path.module}/../../backend/cmd/processor/track-updater/bootstrap.zip"
  source_code_hash = filebase64sha256("${path.module}/../../backend/cmd/processor/track-updater/bootstrap.zip")

  memory_size = 256
  timeout     = 30

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.track_updater]
}

resource "aws_cloudwatch_log_group" "track_updater" {
  name              = "/aws/lambda/${local.name_prefix}-track-updater"
  retention_in_days = 30
}

# Track Creator Lambda
resource "aws_lambda_function" "track_creator" {
  function_name = "${local.name_prefix}-track-creator"
  role          = local.lambda_role_arn
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
      AUDIO_PIPELINE_ARN  = aws_sfn_state_machine.audio_pipeline.arn
    }
  }

  depends_on = [aws_cloudwatch_log_group.track_creator]
}

resource "aws_cloudwatch_log_group" "track_creator" {
  name              = "/aws/lambda/${local.name_prefix}-track-creator"
  retention_in_days = 30
}

# File Mover Lambda
resource "aws_lambda_function" "file_mover" {
  function_name = "${local.name_prefix}-file-mover"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 256
  timeout     = 60

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.file_mover]
}

resource "aws_cloudwatch_log_group" "file_mover" {
  name              = "/aws/lambda/${local.name_prefix}-file-mover"
  retention_in_days = 30
}

# Search Indexer Lambda
resource "aws_lambda_function" "search_indexer" {
  function_name = "${local.name_prefix}-search-indexer"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 256
  timeout     = 30

  environment {
    variables = {
      DYNAMODB_TABLE_NAME       = local.dynamodb_table_name
      NIXIESEARCH_FUNCTION_NAME = aws_lambda_function.nixiesearch.function_name
      EMBEDDING_GATEWAY_URL     = aws_apigatewayv2_api.bedrock_gateway.api_endpoint
      EMBEDDING_API_KEY_SECRET  = aws_secretsmanager_secret.bedrock_gateway_api_key.name
    }
  }

  depends_on = [aws_cloudwatch_log_group.search_indexer]
}

resource "aws_cloudwatch_log_group" "search_indexer" {
  name              = "/aws/lambda/${local.name_prefix}-search-indexer"
  retention_in_days = 30
}

# Upload Status Updater Lambda
resource "aws_lambda_function" "upload_status_updater" {
  function_name = "${local.name_prefix}-upload-status-updater"
  role          = local.lambda_role_arn
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
    }
  }

  depends_on = [aws_cloudwatch_log_group.upload_status_updater]
}

resource "aws_cloudwatch_log_group" "upload_status_updater" {
  name              = "/aws/lambda/${local.name_prefix}-upload-status-updater"
  retention_in_days = 30
}


# Bedrock access policy for audio analyzer
resource "aws_iam_role_policy" "audio_analyzer_bedrock" {
  name = "${local.name_prefix}-audio-analyzer-bedrock"
  role = split("/", local.lambda_role_arn)[1]

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "bedrock:InvokeModel"
        ]
        Resource = [
          "arn:aws:bedrock:${var.aws_region}::foundation-model/anthropic.claude-3-haiku-*",
          "arn:aws:bedrock:${var.aws_region}::foundation-model/anthropic.claude-3-5-haiku-*",
          "arn:aws:bedrock:${var.aws_region}::foundation-model/twelvelabs.marengo-embed-3-0-v1:0"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "bedrock:InvokeModel",
          "bedrock:GetAsyncInvoke"
        ]
        Resource = [
          "arn:aws:bedrock:${var.aws_region}:${data.aws_caller_identity.current.account_id}:async-invoke/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = [
          "arn:aws:s3:::${data.aws_caller_identity.current.account_id}-${local.name_prefix}-media/*"
        ]
      }
    ]
  })
}
