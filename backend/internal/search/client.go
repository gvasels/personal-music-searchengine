package search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaInvoker defines the interface for invoking Lambda functions.
type LambdaInvoker interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

// Client provides search operations via Nixiesearch Lambda.
type Client struct {
	lambdaClient LambdaInvoker
	functionName string
}

// NewClient creates a new search client.
func NewClient(lambdaClient LambdaInvoker, functionName string) *Client {
	return &Client{
		lambdaClient: lambdaClient,
		functionName: functionName,
	}
}

// Search executes a search query and returns results.
func (c *Client) Search(ctx context.Context, userID string, query SearchQuery) (*SearchResponse, error) {
	// Add user filter to scope results
	query.Filters.UserID = userID

	req := NixiesearchRequest{
		Operation: "search",
		Payload:   query,
	}

	resp, err := c.invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	var searchResp SearchResponse
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	if err := json.Unmarshal(data, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	return &searchResp, nil
}

// Index adds or updates a document in the search index.
func (c *Client) Index(ctx context.Context, doc Document) (*IndexResponse, error) {
	req := NixiesearchRequest{
		Operation: "index",
		Payload:   IndexRequest{Document: doc},
	}

	resp, err := c.invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("index failed: %w", err)
	}

	var indexResp IndexResponse
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	if err := json.Unmarshal(data, &indexResp); err != nil {
		return nil, fmt.Errorf("failed to parse index response: %w", err)
	}

	return &indexResp, nil
}

// Delete removes a document from the search index.
func (c *Client) Delete(ctx context.Context, docID string) (*DeleteResponse, error) {
	req := NixiesearchRequest{
		Operation: "delete",
		Payload:   DeleteRequest{ID: docID},
	}

	resp, err := c.invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("delete failed: %w", err)
	}

	var deleteResp DeleteResponse
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	if err := json.Unmarshal(data, &deleteResp); err != nil {
		return nil, fmt.Errorf("failed to parse delete response: %w", err)
	}

	return &deleteResp, nil
}

// BulkIndex adds multiple documents to the search index.
func (c *Client) BulkIndex(ctx context.Context, docs []Document) (*BulkIndexResponse, error) {
	req := NixiesearchRequest{
		Operation: "bulk_index",
		Payload:   BulkIndexRequest{Documents: docs},
	}

	resp, err := c.invoke(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("bulk index failed: %w", err)
	}

	var bulkResp BulkIndexResponse
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	if err := json.Unmarshal(data, &bulkResp); err != nil {
		return nil, fmt.Errorf("failed to parse bulk index response: %w", err)
	}

	return &bulkResp, nil
}

// invoke calls the Nixiesearch Lambda function.
func (c *Client) invoke(ctx context.Context, req NixiesearchRequest) (*NixiesearchResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	result, err := c.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: &c.functionName,
		Payload:      payload,
	})
	if err != nil {
		return nil, fmt.Errorf("lambda invocation failed: %w", err)
	}

	if result.FunctionError != nil {
		return nil, fmt.Errorf("lambda function error: %s", *result.FunctionError)
	}

	var resp NixiesearchResponse
	if err := json.Unmarshal(result.Payload, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse lambda response: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("search operation failed: %s", resp.Error)
	}

	return &resp, nil
}
