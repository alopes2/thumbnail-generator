package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsClient struct {
	s3  s3.Client
	ctx *context.Context
}

func handleRequest(ctx context.Context, event events.S3Event) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("Could not load AWS default configuration")
		return err
	}

	awsClient := awsClient{s3: *s3.NewFromConfig(awsConfig), ctx: &ctx}

	for _, record := range event.Records {
		bucketName := record.S3.Bucket.Name
		objectKey := record.S3.Object.Key

		file, err := awsClient.downloadFile(bucketName, objectKey)

		if err != nil {
			log.Fatalf("Error loading file %s from bucket %s", objectKey, bucketName)
			return err
		}

		err = awsClient.uploadFile(bucketName, objectKey, file)

		if err != nil {
			log.Fatalf("Error uploading file %s to thumbnails/ in bucket %s", objectKey, bucketName)
			return err
		}
	}

	return nil
}

func (client *awsClient) downloadFile(bucketName string, objectKey string) (*bytes.Reader, error) {
	result, err := client.s3.GetObject(*client.ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Fatalf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Fatalf("Error reading file. Error: %s", err)
		return nil, err
	}

	file := bytes.NewReader(body)

	return file, err
}

func (client *awsClient) uploadFile(bucketName string, originalObjectKey string, file *bytes.Reader) error {

	objectKeyParts := strings.Split(originalObjectKey, "/")
	objectKey := fmt.Sprintf("thumbnails/%s", objectKeyParts[len(objectKeyParts)-1])

	_, err := client.s3.PutObject(*client.ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Fatalf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			originalObjectKey, bucketName, objectKey, err)
	}

	return err
}

func main() {
	lambda.Start(handleRequest)
}
