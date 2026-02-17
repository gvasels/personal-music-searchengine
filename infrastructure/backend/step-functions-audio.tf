# Audio Understanding Pipeline Step Function
# Orchestrates: Audio Analyzer -> Embedding Generator -> Track Updater

resource "aws_sfn_state_machine" "audio_pipeline" {
  name     = "${local.name_prefix}-audio-pipeline"
  role_arn = aws_iam_role.audio_pipeline_sfn.arn

  definition = jsonencode({
    Comment = "Audio Understanding Pipeline - analyzes audio and generates embeddings"
    StartAt = "AudioFeatures"
    States = {
      AudioFeatures = {
        Type     = "Task"
        Resource = "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.name_prefix}-audio-features"
        ResultPath = "$.featuresResult"
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.featuresError"
            Next        = "AudioAnalyzer"
          }
        ]
        Next = "AudioAnalyzer"
      }

      AudioAnalyzer = {
        Type     = "Task"
        Resource = aws_lambda_function.audio_analyzer.arn
        ResultPath = "$.analyzerResult"
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.analyzerError"
            Next        = "EmbeddingGenerator"
          }
        ]
        Next = "EmbeddingGenerator"
      }

      EmbeddingGenerator = {
        Type     = "Task"
        Resource = aws_lambda_function.embedding_generator.arn
        ResultPath = "$.embeddingResult"
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.embeddingError"
            Next        = "TrackUpdater"
          }
        ]
        Next = "TrackUpdater"
      }

      TrackUpdater = {
        Type     = "Task"
        Resource = aws_lambda_function.track_updater.arn
        InputPath = "$"
        End       = true
      }
    }
  })

  logging_configuration {
    log_destination        = "${aws_cloudwatch_log_group.audio_pipeline_sfn.arn}:*"
    include_execution_data = true
    level                  = "ERROR"
  }

  tags = {
    Name        = "${local.name_prefix}-audio-pipeline"
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "audio_pipeline_sfn" {
  name              = "/aws/states/${local.name_prefix}-audio-pipeline"
  retention_in_days = 14
}

# IAM Role for Step Function
resource "aws_iam_role" "audio_pipeline_sfn" {
  name = "${local.name_prefix}-audio-pipeline-sfn"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "states.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "audio_pipeline_sfn" {
  name = "lambda-invoke"
  role = aws_iam_role.audio_pipeline_sfn.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = [
          "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.name_prefix}-audio-features",
          aws_lambda_function.audio_analyzer.arn,
          aws_lambda_function.embedding_generator.arn,
          aws_lambda_function.track_updater.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogDelivery",
          "logs:GetLogDelivery",
          "logs:UpdateLogDelivery",
          "logs:DeleteLogDelivery",
          "logs:ListLogDeliveries",
          "logs:PutResourcePolicy",
          "logs:DescribeResourcePolicies",
          "logs:DescribeLogGroups"
        ]
        Resource = "*"
      }
    ]
  })
}

# Output
output "audio_pipeline_arn" {
  value = aws_sfn_state_machine.audio_pipeline.arn
}
