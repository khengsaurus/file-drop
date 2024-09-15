package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/utils"
)

type S3Client struct {
	instance *s3.S3
}

func InitS3Client() *S3Client {
	var awsSession *session.Session
	var err error
	if consts.Local {
		fmt.Println("AWS config: local")
		awsSession, err = session.NewSession(
			&aws.Config{
				Region:   aws.String(os.Getenv("AWS_REGION")),
				Endpoint: aws.String(os.Getenv("AWS_URI_DEV")),
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY_DEV"),
					os.Getenv("AWS_SECRET_KEY_DEV"), ""),
				S3ForcePathStyle: aws.Bool(true),
			})
	} else {
		fmt.Println("AWS config: remote")
		awsSession, err = session.NewSession(
			&aws.Config{
				Region: aws.String(os.Getenv("AWS_REGION")),
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("AWS_ACCESS_KEY"),
					os.Getenv("AWS_SECRET_KEY"),
					"",
				),
			})
	}
	if err != nil {
		log.Fatal(err)
	}

	return &S3Client{instance: s3.New(awsSession)}
}

func GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	s3Client, ok := ctx.Value(consts.S3ClientKey).(*S3Client)
	if !ok {
		return nil, fmt.Errorf("couldn't find %s in request context", consts.S3ClientKey)
	}

	result, err := s3Client.instance.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key)},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// finds a key that does not exist in either redis or s3
func GetSignedPutUrl(
	ctx context.Context,
	contentType string,
	size int,
	attempt int,
	clientToken string,
) (string, string, error) {
	if attempt > consts.WriteTries {
		return "", "", fmt.Errorf("failed to get presignedPutUrl after %d tries", consts.WriteTries)
	}

	s3Client, ok := ctx.Value(consts.S3ClientKey).(*S3Client)
	if !ok {
		return "", "", fmt.Errorf("couldn't find %s in request context", consts.S3ClientKey)
	}

	redisClient, err := GetRedisClient(ctx)
	if err != nil {
		return "", "", err
	}

	key := utils.RandString(6)
	_, err = s3Client.instance.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key),
	})
	if err == nil || !strings.Contains(err.Error(), "NoSuchKey") {
		// no error - object with key exists in s3
		return GetSignedPutUrl(ctx, contentType, size, attempt+1, clientToken)
	}

	exists, _ := redisClient.CheckExists(ctx, key)
	if exists {
		return GetSignedPutUrl(ctx, contentType, size, attempt+1, clientToken)
	}

	placeholder := fmt.Sprintf("%s_%s", consts.RedisValPlaceholderPrefix, clientToken)
	err = redisClient.SetValue(ctx, key, placeholder, time.Hour)
	if err != nil {
		return "", "", err
	}

	req, _ := s3Client.instance.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket:        aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:           aws.String(key),
			ContentType:   aws.String(contentType),
			ContentLength: aws.Int64(int64(size)),
		},
	)
	url, err := req.Presign(1 * time.Minute)
	return key, url, err
}

func GetSignedGetUrl(ctx context.Context, key string) (string, error) {
	s3Client, ok := ctx.Value(consts.S3ClientKey).(*S3Client)
	if !ok {
		return "", fmt.Errorf("couldn't find %s in request context", consts.S3ClientKey)
	}

	req, _ := s3Client.instance.GetObjectRequest(
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:    aws.String(key),
		},
	)

	return req.Presign(24 * time.Hour)
}

func DeleteObject(ctx context.Context, key string) (bool, error) {
	s3Client, ok := ctx.Value(consts.S3ClientKey).(*S3Client)
	if !ok {
		return false, fmt.Errorf("couldn't find %s in request context", consts.S3ClientKey)
	}

	if _, err := s3Client.instance.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
			Key:    aws.String(key),
		},
	); err != nil {
		return false, err
	}

	return true, nil
}

func DeleteObjects(
	ctx context.Context,
	objects *s3.DeleteObjectsInput,
) (bool, error) {
	s3Client, ok := ctx.Value(consts.S3ClientKey).(*S3Client)
	if !ok {
		return false, fmt.Errorf("couldn't find %s in request context", consts.S3ClientKey)
	}

	_, err := s3Client.instance.DeleteObjects(objects)
	if err != nil {
		return false, err
	}

	return true, err
}
