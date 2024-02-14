package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsClient struct {
	s3 s3.Client
}

func handleRequest(ctx context.Context, event events.SQSEvent) error {
	log.Printf("Received event records %v", event.Records)
	log.Printf("Received event first record %v", event.Records[0].Body)

	// awsConfig, err := config.LoadDefaultConfig(ctx)

	// if err != nil {
	// 	log.Fatalf("Could not load AWS default configuration")
	// 	return err
	// }

	// awsClient := awsClient{s3: *s3.NewFromConfig(awsConfig)}

	// for _, record := range event.Records {
	// 	log.Printf("Processing SQS record %v", record)
	// 	var s3EventRecord events.S3EventRecord

	// 	err := json.Unmarshal([]byte(record.Body), &s3EventRecord)

	// 	log.Printf("Unmarshalled s3 event %v", s3EventRecord)
	// 	if err != nil {
	// 		log.Fatalf("Could not unmarshal SQS Body %s to S3 Event Record", record.Body)
	// 		return err
	// 	}

	// 	log.Printf("S3 %v", s3EventRecord.S3)
	// 	bucketName := s3EventRecord.S3.Bucket.Name
	// 	objectKey := s3EventRecord.S3.Object.Key

	// 	file, err := awsClient.downloadFile(bucketName, objectKey)

	// 	if err != nil {
	// 		log.Fatalf("Error loading file %s from bucket %s", objectKey, bucketName)
	// 		return err
	// 	}

	// 	err = awsClient.uploadFile(bucketName, objectKey, file)

	// 	if err != nil {
	// 		log.Fatalf("Error uploading file %s to thumbnails/ in bucket %s", objectKey, bucketName)
	// 		return err
	// 	}
	// }

	return nil
}

func (client *awsClient) downloadFile(bucketName string, objectKey string) (*os.File, error) {
	fileName := ""

	result, err := client.s3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return nil, err
	}

	defer result.Body.Close()

	file, err := os.Create(fileName)

	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return nil, err
	}

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}

	_, err = file.Write(body)

	return file, err
}

func (client *awsClient) uploadFile(bucketName string, originalObjectKey string, file *os.File) error {

	defer file.Close()

	objectKeyParts := strings.Split(originalObjectKey, "/")
	objectKey := fmt.Sprintf("thumbnails/%s", objectKeyParts[len(objectKeyParts)-1])

	_, err := client.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			originalObjectKey, bucketName, objectKey, err)
	}

	return err
}

func main() {
	lambda.Start(handleRequest)
}
