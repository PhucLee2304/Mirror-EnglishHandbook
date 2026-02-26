package minio

import (
	"context"
	"core/config"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func EnsureBucketPublic(client *minio.Client, bucketName string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking if bucket exists %s: %w", bucketName, err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket at %s: %w", bucketName, err)
		}
		log.Printf("bucket created at %s", bucketName)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetBucketLocation", "s3:ListBucket"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s"]
			},
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName, bucketName)
	err = client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return fmt.Errorf("error setting bucket policy: %w", err)
	}
	log.Printf("bucket policy set at %s", bucketName)
	return nil
}

func InitMinioClient(cfg *config.Config) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.MinioHost, cfg.MinioPort)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating minio client: %w", err)
	}

	err = EnsureBucketPublic(client, cfg.BucketBooksVideoName)
	if err != nil {
		return nil, fmt.Errorf("error checking if bucket books video exists: %w", err)
	}
	return client, nil
}
