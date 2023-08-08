package clients

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const PrimaryBucketName = "rpg-vault-go"

type S3ClientWrapper struct {
	*s3.Client
}

func (c *S3ClientWrapper) Upload(contents []byte, destinationPath string) (string, error) {
	contentsReader := bytes.NewReader(contents)
	_, err := c.PutObject(context.TODO(), &s3.PutObjectInput{
		//Bucket: aws.String(config.PrimaryBucketName),
		Bucket: aws.String(PrimaryBucketName),
		Key:    aws.String(destinationPath),
		Body:   contentsReader,
	})
	if err != nil {
		return "", err
	}

	//https://my-bucket.s3-ap-southeast-2.amazonaws.com/foo/bar.txt
	//https://rpg-vault-go.s3.sa-east-1.amazonaws.com/f572a37c-33a7-4b5b-85d9-86cb596d2edb/world/cover.png
	//TODO: get from config
	region := "sa-east-1"

	//TODO: get from config?
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", PrimaryBucketName, region, destinationPath), nil
}

func (c *S3ClientWrapper) Read(path string) []byte {
	output, err := c.GetObject(context.TODO(), &s3.GetObjectInput{
		//Bucket: aws.String(config.PrimaryBucketName),
		Bucket: aws.String(PrimaryBucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(output.Body)
	return buf.Bytes()
}

func (c *S3ClientWrapper) Delete(path string) {
	_, err := c.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		//Bucket: aws.String(config.PrimaryBucketName),
		Bucket: aws.String(PrimaryBucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		panic(err)
	}
}
