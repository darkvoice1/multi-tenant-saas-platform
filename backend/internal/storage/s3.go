package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/darkvoice1/multi-tenant-saas-platform/backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3(cfg config.Config) (*S3Storage, error) {
	if cfg.S3Endpoint == "" || cfg.S3AccessKey == "" || cfg.S3SecretKey == "" {
		return nil, fmt.Errorf("s3 config missing")
	}
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && cfg.S3Endpoint != "" {
			return aws.Endpoint{
				URL:               cfg.S3Endpoint,
				SigningRegion:     cfg.S3Region,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(endpointResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("load s3 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.S3PathStyle
	})

	st := &S3Storage{
		client: client,
		bucket: cfg.S3Bucket,
	}
	if err := st.ensureBucket(context.Background(), cfg); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *S3Storage) ensureBucket(ctx context.Context, cfg config.Config) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(s.bucket)})
	if err == nil {
		return nil
	}
	endpointHost := ""
	if cfg.S3Endpoint != "" {
		if u, parseErr := url.Parse(cfg.S3Endpoint); parseErr == nil {
			endpointHost = u.Hostname()
		}
	}
	input := &s3.CreateBucketInput{Bucket: aws.String(s.bucket)}
	if endpointHost == "" && cfg.S3Region != "" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(cfg.S3Region),
		}
	}
	_, err = s.client.CreateBucket(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("create bucket: %w", err)
		}
		code := apiErr.ErrorCode()
		if code != "BucketAlreadyOwnedByYou" && code != "BucketAlreadyExists" {
			return fmt.Errorf("create bucket: %w", err)
		}
	}
	return nil
}

func (s *S3Storage) Save(ctx context.Context, key string, body io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: &size,
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *S3Storage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return obj.Body, nil
}
