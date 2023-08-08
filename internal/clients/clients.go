package clients

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfigMod "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/aws/aws-sdk-go-v2/service/ses"
	// "github.com/aws/aws-sdk-go-v2/service/sns"
)

var awsConfig aws.Config
var onceAwsConfig sync.Once

func getAwsConfig(awsKey, awsSecret string) aws.Config {
	onceAwsConfig.Do(func() {
		var err error

		credentialProvider := awsConfigMod.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsKey, awsSecret, ""))

		awsConfig, err = awsConfigMod.LoadDefaultConfig(context.TODO(), credentialProvider)
		if err != nil {
			panic(err)
		}
	})

	return awsConfig
}

func GetDynamodbClient(awsKey, awsSecret, awsRegion string) *DynamoDbClientWrapper {
	awsConfig = getAwsConfig(awsKey, awsSecret)

	dynamodbClient := dynamodb.NewFromConfig(awsConfig, func(opt *dynamodb.Options) {
		opt.Region = awsRegion
	})

	return &DynamoDbClientWrapper{
		Client: dynamodbClient,
	}
}

func GetS3Client(awsKey, awsSecret, awsRegion string) *S3ClientWrapper {
	awsConfig = getAwsConfig(awsKey, awsSecret)

	s3Client := s3.NewFromConfig(awsConfig, func(opt *s3.Options) {
		opt.Region = awsRegion
	})

	return &S3ClientWrapper{
		Client: s3Client,
	}
}
