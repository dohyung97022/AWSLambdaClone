package lambdaClone

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type S3Client struct {
	client *s3.Client
}

var s3Client S3Client

const bucketName string = "aws-lambda-clone"
const region string = "ap-northeast-2"

func getEC2CredentialsConfig() (aws.Config, error) {
	// validate if ec2 credentials exist
	_, err := ec2rolecreds.New().Retrieve(context.Background())
	if err != nil {
		return aws.Config{}, err
	}
	return config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(ec2rolecreds.New()),
		config.WithRegion(region),
	)
}

func getDefaultCredentialsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)
}

func SetS3CClient() error {
	cfg, err := getEC2CredentialsConfig()
	if err != nil {
		cfg, err = getDefaultCredentialsConfig()
	}
	if err != nil {
		return err
	}
	s3Client = S3Client{client: s3.NewFromConfig(cfg)}
	return nil
}

func (s3Client S3Client) uploadCodeToS3(lambda *Lambda) error {
	_, err := s3Client.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(lambda.Id),
		Body:   bytes.NewReader([]byte(lambda.Code)),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3Client S3Client) readCodeFromS3(lambda *Lambda) (string, error) {
	result, err := s3Client.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(lambda.Id),
	})
	if err != nil {
		return "", err
	}
	var output = make([]byte, *result.ContentLength)
	if _, err = result.Body.Read(output); err != nil && err != io.EOF {
		return "", err
	}
	if err = result.Body.Close(); err != nil {
		return "", err
	}
	return string(output), nil
}
