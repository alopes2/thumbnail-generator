package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

type awsClient struct {
	s3  s3.Client
	ctx *context.Context
}

func handleRequest(ctx context.Context, event events.SNSEvent) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatalf("Could not load AWS default configuration")
		return err
	}

	awsClient := awsClient{s3: *s3.NewFromConfig(awsConfig), ctx: &ctx}

	for _, record := range event.Records {
		var imageEvent events.S3Event

		err := json.Unmarshal([]byte(record.SNS.Message), &imageEvent)

		if err != nil {
			log.Fatalf("Could not unmarshal SNS message %s into S3 Event Record with error: %v", record.SNS.Message, err)
			return err
		}

		for _, imageRecord := range imageEvent.Records {
			bucketName := imageRecord.S3.Bucket.Name
			objectKey := imageRecord.S3.Object.Key

			file, err := awsClient.downloadFile(bucketName, objectKey)

			log.Printf("Successfully downloaded image")

			if err != nil {
				log.Fatalf("Error loading file %s from bucket %s", objectKey, bucketName)
				return err
			}

			thumbnail, err := createThumbnail(file)

			if err != nil {
				log.Fatalf("Error creating thumbnail for file %s from bucket %s. Error is %v", objectKey, bucketName, err)
				return err
			}

			log.Printf("Successfully created thumbnail")

			err = awsClient.uploadFile(bucketName, objectKey, thumbnail)

			log.Printf("Successfully uploaded thumbnail")

			if err != nil {
				log.Fatalf("Error uploading file %s to thumbnails/ in bucket %s", objectKey, bucketName)
				return err
			}
		}
	}

	return nil
}

func createThumbnail(reader io.Reader) (*bytes.Buffer, error) {
	srcImage, _, err := image.Decode(reader)

	if err != nil {
		log.Fatalf("Could not decode file because of error %v", err)
		return nil, err
	}

	// Generates a 80x80 thumbnail
	thumbnail := imaging.Thumbnail(srcImage, 80, 80, imaging.Lanczos)

	var bufferBytes []byte
	buffer := bytes.NewBuffer(bufferBytes)

	err = png.Encode(buffer, thumbnail)

	return buffer, err
}

func (client *awsClient) downloadFile(bucketName string, objectKey string) (*bytes.Reader, error) {
	result, err := client.s3.GetObject(*client.ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Fatalf("Couldn't get object %v:%v. Here's why: %v", bucketName, objectKey, err)
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

func (client *awsClient) uploadFile(bucketName string, originalObjectKey string, thumbnail io.Reader) error {

	objectKeyParts := strings.Split(originalObjectKey, "/")
	fileNameWithoutExtensions := strings.Split(objectKeyParts[len(objectKeyParts)-1], ".")[0]
	objectKey := fmt.Sprintf("thumbnails/%s_thumbnail.png", fileNameWithoutExtensions)

	_, err := client.s3.PutObject(*client.ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   thumbnail,
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
