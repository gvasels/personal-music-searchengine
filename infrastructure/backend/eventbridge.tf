# EventBridge Rules for MediaConvert and Scheduled Tasks

# EventBridge Rule for MediaConvert Job Completion
resource "aws_cloudwatch_event_rule" "mediaconvert_complete" {
  name        = "${local.name_prefix}-mediaconvert-complete"
  description = "Trigger Lambda when MediaConvert job completes"

  event_pattern = jsonencode({
    source      = ["aws.mediaconvert"]
    detail-type = ["MediaConvert Job State Change"]
    detail = {
      status = ["COMPLETE"]
      queue  = [aws_media_convert_queue.default.arn]
    }
  })
}

# EventBridge Rule for MediaConvert Job Failure
resource "aws_cloudwatch_event_rule" "mediaconvert_error" {
  name        = "${local.name_prefix}-mediaconvert-error"
  description = "Trigger Lambda when MediaConvert job fails"

  event_pattern = jsonencode({
    source      = ["aws.mediaconvert"]
    detail-type = ["MediaConvert Job State Change"]
    detail = {
      status = ["ERROR", "CANCELED"]
      queue  = [aws_media_convert_queue.default.arn]
    }
  })
}

# Target for MediaConvert completion - invoke transcode complete Lambda
resource "aws_cloudwatch_event_target" "mediaconvert_complete" {
  rule      = aws_cloudwatch_event_rule.mediaconvert_complete.name
  target_id = "TranscodeComplete"
  arn       = aws_lambda_function.transcode_complete.arn
}

# Target for MediaConvert error - invoke transcode complete Lambda (handles failures)
resource "aws_cloudwatch_event_target" "mediaconvert_error" {
  rule      = aws_cloudwatch_event_rule.mediaconvert_error.name
  target_id = "TranscodeError"
  arn       = aws_lambda_function.transcode_complete.arn
}

# Permission for EventBridge to invoke transcode complete Lambda (success)
resource "aws_lambda_permission" "eventbridge_transcode_complete" {
  statement_id  = "AllowEventBridgeComplete"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.transcode_complete.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.mediaconvert_complete.arn
}

# Permission for EventBridge to invoke transcode complete Lambda (error)
resource "aws_lambda_permission" "eventbridge_transcode_error" {
  statement_id  = "AllowEventBridgeError"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.transcode_complete.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.mediaconvert_error.arn
}

# Scheduled Rule for Daily Index Rebuild
resource "aws_cloudwatch_event_rule" "daily_index_rebuild" {
  name                = "${local.name_prefix}-daily-index-rebuild"
  description         = "Trigger daily search index rebuild"
  schedule_expression = "cron(0 3 * * ? *)" # 3 AM UTC daily
}

# Index Rebuild Lambda (reuses search indexer with rebuild mode)
resource "aws_lambda_function" "index_rebuild" {
  function_name = "${local.name_prefix}-index-rebuild"
  role          = local.lambda_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]

  filename         = data.archive_file.placeholder.output_path
  source_code_hash = data.archive_file.placeholder.output_base64sha256

  memory_size = 512
  timeout     = 300 # 5 minutes for full rebuild

  environment {
    variables = {
      DYNAMODB_TABLE_NAME       = local.dynamodb_table_name
      NIXIESEARCH_FUNCTION_NAME = aws_lambda_function.nixiesearch.function_name
      REBUILD_MODE              = "true"
    }
  }

  depends_on = [aws_cloudwatch_log_group.index_rebuild]
}

resource "aws_cloudwatch_log_group" "index_rebuild" {
  name              = "/aws/lambda/${local.name_prefix}-index-rebuild"
  retention_in_days = 30
}

# Target for daily index rebuild
resource "aws_cloudwatch_event_target" "daily_index_rebuild" {
  rule      = aws_cloudwatch_event_rule.daily_index_rebuild.name
  target_id = "IndexRebuild"
  arn       = aws_lambda_function.index_rebuild.arn
}

# Permission for EventBridge to invoke index rebuild Lambda
resource "aws_lambda_permission" "eventbridge_index_rebuild" {
  statement_id  = "AllowEventBridgeRebuild"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.index_rebuild.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily_index_rebuild.arn
}

# Outputs
output "mediaconvert_complete_rule_arn" {
  description = "EventBridge rule ARN for MediaConvert completion"
  value       = aws_cloudwatch_event_rule.mediaconvert_complete.arn
}

output "daily_rebuild_rule_arn" {
  description = "EventBridge rule ARN for daily index rebuild"
  value       = aws_cloudwatch_event_rule.daily_index_rebuild.arn
}

output "index_rebuild_lambda_arn" {
  description = "Index rebuild Lambda ARN"
  value       = aws_lambda_function.index_rebuild.arn
}
