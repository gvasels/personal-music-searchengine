package clients

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManager implements SecretsManagerClient using the AWS SDK.
type AWSSecretsManager struct {
	client *secretsmanager.Client
}

// NewAWSSecretsManager creates a new AWS Secrets Manager client adapter.
func NewAWSSecretsManager(client *secretsmanager.Client) *AWSSecretsManager {
	return &AWSSecretsManager{client: client}
}

// GetSecretString fetches a secret string value by name.
func (s *AWSSecretsManager) GetSecretString(ctx context.Context, secretName string) (string, error) {
	out, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}
	if out.SecretString == nil {
		return "", fmt.Errorf("secret %s has no string value", secretName)
	}
	return *out.SecretString, nil
}
