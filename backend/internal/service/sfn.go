package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

// SFNClientAdapter adapts the AWS SFN client to our StepFunctionsClient interface
type SFNClientAdapter struct {
	client *sfn.Client
}

// NewSFNClientAdapter creates a new Step Functions client adapter
func NewSFNClientAdapter(client *sfn.Client) *SFNClientAdapter {
	return &SFNClientAdapter{client: client}
}

// StartExecution starts a Step Functions execution
func (a *SFNClientAdapter) StartExecution(ctx context.Context, input *StepFunctionsStartInput) (*StepFunctionsStartOutput, error) {
	result, err := a.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(input.StateMachineArn),
		Name:            aws.String(input.Name),
		Input:           aws.String(input.Input),
	})
	if err != nil {
		return nil, err
	}

	return &StepFunctionsStartOutput{
		ExecutionArn: aws.ToString(result.ExecutionArn),
		StartDate:    aws.ToTime(result.StartDate),
	}, nil
}
