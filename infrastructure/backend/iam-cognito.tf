# Cognito Admin IAM Policy for Admin User Management
# Task 4.2: Admin panel requires Cognito user management permissions

# Cognito Admin Policy - allows Lambda to manage Cognito users
resource "aws_iam_policy" "cognito_admin" {
  name        = "${local.name_prefix}-cognito-admin"
  description = "Allow Lambda to manage Cognito users for admin panel"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "CognitoUserManagement"
        Effect = "Allow"
        Action = [
          "cognito-idp:AdminGetUser",
          "cognito-idp:AdminListGroupsForUser",
          "cognito-idp:AdminAddUserToGroup",
          "cognito-idp:AdminRemoveUserFromGroup",
          "cognito-idp:AdminDisableUser",
          "cognito-idp:AdminEnableUser",
          "cognito-idp:ListUsers",
          "cognito-idp:ListUsersInGroup"
        ]
        Resource = local.cognito_user_pool_arn
      }
    ]
  })
}

# Attach Cognito admin policy to Lambda execution role
resource "aws_iam_role_policy_attachment" "lambda_cognito_admin" {
  role       = local.lambda_role_name
  policy_arn = aws_iam_policy.cognito_admin.arn
}
