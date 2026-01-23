# Step Functions Processor Lambda Functions

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
resource "aws_lambda_layer_version" "ffmpeg" {
  layer_name          = "${local.name_prefix}-ffmpeg"
  filename            = "${path.module}/ffmpeg-layer.zip"
  source_code_hash    = filebase64sha256("${path.module}/ffmpeg-layer.zip")
  compatible_runtimes = ["provided.al2023"]
  compatible_architectures = ["arm64"]
  description         = "FFmpeg static binaries for ARM64"
}

# Audio Analyzer Lambda (BPM and Key Detection)
resource "aws_lambda_function" "audio_analyzer" {
  function_name = "${local.name_prefix}-audio-analyzer"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 1024  # Audio analysis may need more memory
  timeout     = 60    # Allow more time for FFmpeg processing

  # Custom FFmpeg Lambda layer for audio processing (ARM64)
  layers = [aws_lambda_layer_version.ffmpeg.arn]

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = local.dynamodb_table_name
      MEDIA_BUCKET        = local.media_bucket_name
      FFMPEG_PATH         = "/opt/bin/ffmpeg"
      FFPROBE_PATH        = "/opt/bin/ffprobe"
    }
  }

  depends_on = [aws_cloudwatch_log_group.audio_analyzer]
}

resource "aws_cloudwatch_log_group" "audio_analyzer" {
  name              = "/aws/lambda/${local.name_prefix}-audio-analyzer"
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
