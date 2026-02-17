# IAM Policy for shared Lambda role to read Bedrock Gateway API key
resource "aws_iam_role_policy" "lambda_embedding_gateway_secrets" {
  name = "embedding-gateway-secrets"
  role = local.lambda_role_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["secretsmanager:GetSecretValue"]
        Resource = aws_secretsmanager_secret.bedrock_gateway_api_key.arn
      }
    ]
  })
}
