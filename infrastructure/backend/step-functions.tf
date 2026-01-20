# Step Functions State Machine for Upload Processing
# Orchestrates: Metadata Extraction → Cover Art Processing → Search Indexing

resource "aws_sfn_state_machine" "upload_processor" {
  name     = "${local.name_prefix}-upload-processor"
  role_arn = aws_iam_role.step_functions.arn

  definition = jsonencode({
    Comment = "Process uploaded audio files: extract metadata, process cover art, index for search"
    StartAt = "ExtractMetadata"
    States = {
      ExtractMetadata = {
        Type     = "Task"
        Resource = aws_lambda_function.metadata_extractor.arn
        Parameters = {
          "uploadId.$"  = "$.uploadId"
          "userId.$"    = "$.userId"
          "s3Key.$"     = "$.s3Key"
          "fileName.$"  = "$.fileName"
          "bucketName"  = local.media_bucket_name
        }
        ResultPath = "$.metadata"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.error"
            Next        = "MarkUploadFailed"
          }
        ]
        Next = "ProcessCoverArt"
      }

      ProcessCoverArt = {
        Type     = "Task"
        Resource = aws_lambda_function.cover_art_processor.arn
        Parameters = {
          "uploadId.$"    = "$.uploadId"
          "userId.$"      = "$.userId"
          "s3Key.$"       = "$.s3Key"
          "metadata.$"    = "$.metadata"
          "bucketName"    = local.media_bucket_name
        }
        ResultPath = "$.coverArt"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.coverArtError"
            Next        = "CreateTrackRecord"  # Continue even if cover art fails
          }
        ]
        Next = "CreateTrackRecord"
      }

      CreateTrackRecord = {
        Type     = "Task"
        Resource = aws_lambda_function.track_creator.arn
        Parameters = {
          "uploadId.$"   = "$.uploadId"
          "userId.$"     = "$.userId"
          "s3Key.$"      = "$.s3Key"
          "fileName.$"   = "$.fileName"
          "metadata.$"   = "$.metadata"
          "coverArt.$"   = "$.coverArt"
          "bucketName"   = local.media_bucket_name
          "tableName"    = local.dynamodb_table_name
        }
        ResultPath = "$.track"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.error"
            Next        = "MarkUploadFailed"
          }
        ]
        Next = "MoveToMediaStorage"
      }

      MoveToMediaStorage = {
        Type     = "Task"
        Resource = aws_lambda_function.file_mover.arn
        Parameters = {
          "uploadId.$"   = "$.uploadId"
          "userId.$"     = "$.userId"
          "sourceKey.$"  = "$.s3Key"
          "trackId.$"    = "$.track.trackId"
          "bucketName"   = local.media_bucket_name
        }
        ResultPath = "$.finalLocation"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.error"
            Next        = "MarkUploadFailed"
          }
        ]
        Next = "StartTranscode"
      }

      StartTranscode = {
        Type     = "Task"
        Resource = aws_lambda_function.transcode_start.arn
        Parameters = {
          "trackId.$"      = "$.track.trackId"
          "userId.$"       = "$.userId"
          "s3Key.$"        = "$.finalLocation.newKey"
          "format.$"       = "$.metadata.format"
          "bucketName"     = local.media_bucket_name
          "tableName"      = local.dynamodb_table_name
        }
        ResultPath = "$.transcode"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.transcodeError"
            Next        = "IndexForSearch"  # Continue even if transcode fails to start
          }
        ]
        Next = "IndexForSearch"
      }

      IndexForSearch = {
        Type     = "Task"
        Resource = aws_lambda_function.search_indexer.arn
        Parameters = {
          "trackId.$"   = "$.track.trackId"
          "userId.$"    = "$.userId"
          "metadata.$"  = "$.metadata"
          "tableName"   = local.dynamodb_table_name
        }
        ResultPath = "$.searchIndex"
        Retry = [
          {
            ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException"]
            IntervalSeconds = 2
            MaxAttempts     = 3
            BackoffRate     = 2
          }
        ]
        Catch = [
          {
            ErrorEquals = ["States.ALL"]
            ResultPath  = "$.searchError"
            Next        = "MarkUploadCompleted"  # Continue even if indexing fails
          }
        ]
        Next = "MarkUploadCompleted"
      }

      MarkUploadCompleted = {
        Type     = "Task"
        Resource = aws_lambda_function.upload_status_updater.arn
        Parameters = {
          "uploadId.$"  = "$.uploadId"
          "userId.$"    = "$.userId"
          "trackId.$"   = "$.track.trackId"
          "status"      = "COMPLETED"
          "tableName"   = local.dynamodb_table_name
        }
        End = true
      }

      MarkUploadFailed = {
        Type     = "Task"
        Resource = aws_lambda_function.upload_status_updater.arn
        Parameters = {
          "uploadId.$"  = "$.uploadId"
          "userId.$"    = "$.userId"
          "status"      = "FAILED"
          "error.$"     = "$.error"
          "tableName"   = local.dynamodb_table_name
        }
        End = true
      }
    }
  })

  logging_configuration {
    log_destination        = "${aws_cloudwatch_log_group.step_functions.arn}:*"
    include_execution_data = true
    level                  = "ERROR"
  }
}

# CloudWatch Log Group for Step Functions
resource "aws_cloudwatch_log_group" "step_functions" {
  name              = "/aws/vendedlogs/states/${local.name_prefix}-upload-processor"
  retention_in_days = 30
}

# IAM Role for Step Functions
resource "aws_iam_role" "step_functions" {
  name = "${local.name_prefix}-step-functions"

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

# IAM Policy for Step Functions to invoke Lambdas
resource "aws_iam_role_policy" "step_functions_lambda" {
  name = "${local.name_prefix}-step-functions-lambda"
  role = aws_iam_role.step_functions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = [
          aws_lambda_function.metadata_extractor.arn,
          aws_lambda_function.cover_art_processor.arn,
          aws_lambda_function.track_creator.arn,
          aws_lambda_function.file_mover.arn,
          aws_lambda_function.transcode_start.arn,
          aws_lambda_function.search_indexer.arn,
          aws_lambda_function.upload_status_updater.arn
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

# IAM Policy for API Lambda to start Step Functions execution
resource "aws_iam_role_policy" "lambda_step_functions" {
  name = "${local.name_prefix}-lambda-step-functions"
  role = split("/", local.lambda_role_arn)[1]  # Extract role name from ARN

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "states:StartExecution",
          "states:DescribeExecution"
        ]
        Resource = aws_sfn_state_machine.upload_processor.arn
      }
    ]
  })
}
