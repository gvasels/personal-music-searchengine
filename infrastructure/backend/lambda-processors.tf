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
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 60

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
      BEDROCK_MODEL_ID    = "global.anthropic.claude-sonnet-4-6"
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
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 120

  environment {
    variables = {
      MEDIA_BUCKET        = local.media_bucket_name
      VECTOR_BUCKET_NAME  = local.vector_bucket_name
      VECTOR_INDEX_NAME   = "media-embeddings"
      AWS_ACCOUNT_ID      = data.aws_caller_identity.current.account_id
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
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
      VECTOR_BUCKET_NAME        = local.vector_bucket_name
      SEARCH_VECTOR_INDEX_NAME  = "search-embeddings"
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


# Audio Features Lambda (Python/librosa for BPM and key detection - container image)
resource "aws_lambda_function" "audio_features" {
  function_name = "${local.name_prefix}-audio-features"
  role          = local.lambda_role_arn
  package_type  = "Image"
  image_uri     = "${data.terraform_remote_state.global.outputs.ecr_repository_urls.audio_features}:latest"
  architectures = ["arm64"]

  memory_size = 1024
  timeout     = 120

  environment {
    variables = {
      MEDIA_BUCKET = local.media_bucket_name
    }
  }

  depends_on = [aws_cloudwatch_log_group.audio_features]
}

resource "aws_cloudwatch_log_group" "audio_features" {
  name              = "/aws/lambda/${local.name_prefix}-audio-features"
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
          "arn:aws:bedrock:${var.aws_region}:${data.aws_caller_identity.current.account_id}:inference-profile/global.anthropic.claude-sonnet-4-6",
          "arn:aws:bedrock:*::foundation-model/anthropic.claude-sonnet-4-6",
          "arn:aws:bedrock:${var.aws_region}:${data.aws_caller_identity.current.account_id}:inference-profile/us.twelvelabs.marengo-embed-3-0-v1:0",
          # Cross-region inference profile routes Marengo to us-east-2 — wildcard region required
          "arn:aws:bedrock:*::foundation-model/us.twelvelabs.marengo-embed-3-0-v1:0",
          "arn:aws:bedrock:*::foundation-model/twelvelabs.marengo-embed-3-0-v1:0"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "bedrock:InvokeModel",
          "bedrock:StartAsyncInvoke",
          "bedrock:GetAsyncInvoke"
        ]
        Resource = [
          "arn:aws:bedrock:${var.aws_region}:${data.aws_caller_identity.current.account_id}:async-invoke/*",
          "arn:aws:bedrock:${var.aws_region}:${data.aws_caller_identity.current.account_id}:inference-profile/us.twelvelabs.marengo-embed-3-0-v1:0",
          # Cross-region inference profile routes Marengo to us-east-2 — wildcard region required
          "arn:aws:bedrock:*::foundation-model/us.twelvelabs.marengo-embed-3-0-v1:0",
          "arn:aws:bedrock:*::foundation-model/twelvelabs.marengo-embed-3-0-v1:0"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject"
        ]
        Resource = [
          "arn:aws:s3:::${local.media_bucket_name}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "aws-marketplace:ViewSubscriptions",
          "aws-marketplace:Subscribe"
        ]
        Resource = "*"
      }
    ]
  })
}
