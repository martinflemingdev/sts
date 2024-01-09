package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// NewAwsConfig creates and returns a new AWS configuration
// This uses the Lambda's execution role and the Lambda's deployed Region by default
func NewAwsConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("error loading AWS configuration: %v", err)
	}
	return cfg, nil
}

// NewS3Client creates a new S3 client using the provided AWS configuration
func NewS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

// AssumeRole creates a new AWS configuration using the specified role ARN
func AssumeRole(ctx context.Context, cfg aws.Config, roleArn string) (aws.Config, error) {
	stsClient := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)
	roleConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(creds))
	if err != nil {
		return aws.Config{}, err
	}
	return roleConfig, nil
}

// NewDynamoDBClient creates a new DynamoDB client using the specified AWS configuration
func NewDynamoDBClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func main() {
	ctx := context.Background()

	// Create an AWS configuration for the Lambda's role
	cfg, err := NewAwsConfig(ctx)
	if err != nil {
		fmt.Println("Error creating AWS configuration for Lambda's role:", err)
		return
	}

	// Create an S3 client using the Lambda's role
	s3Client := NewS3Client(cfg)
	fmt.Println("S3 client created successfully:", s3Client)

	// Assume the DynamoDB read role
	dynamoReadRoleArn := "arn:aws:iam::123456789012:role/DynamoReadRole" // replace with your actual role ARN
	assumedRoleConfig, err := AssumeRole(ctx, cfg, dynamoReadRoleArn)
	if err != nil {
		fmt.Println("Error assuming DynamoDB read role:", err)
		return
	}

	// Create a DynamoDB client using the assumed role's configuration
	dynamoClient := NewDynamoDBClient(assumedRoleConfig)
	fmt.Println("DynamoDB client created successfully:", dynamoClient)

	// Here you can use s3Client and dynamoClient as needed
}

// Key differences in this version:

// Configuration Management: AWS SDK for Go v2 uses a unified aws.Config object instead of sessions. 
// The config.LoadDefaultConfig function is used to load the AWS configuration.

// Client Creation: Clients are created directly from the configuration object using functions like NewFromConfig.

// Context Support: The SDK v2 methods typically require a context. 
// This context is used for timeout and cancellation signals across API calls.

// Assume Role: The stscreds package is used for creating a credentials provider that assumes a role.

// This updated code is aligned with the patterns and practices of AWS SDK for Go v2. 
// Ensure that your environment is set up with the correct version of the AWS SDK for Go v2 and that your 
// AWS credentials are configured properly for this to work as expected.