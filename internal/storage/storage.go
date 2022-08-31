package storage

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	BucketName string
	Key        string
	CodePath   string
	client     *s3.Client
}

func (s *Storage) Initialize(cfg aws.Config) {
	s.client = s3.NewFromConfig(cfg)
}

func (s *Storage) UpdateCode(ctx context.Context) error {
	log.Printf("opening %s\n", s.CodePath)

	dat, err := os.ReadFile(s.CodePath)
	if err != nil {
		return err
	}

	if _, err = s.client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(s.Key),
			Body:   bytes.NewReader(dat),
		},
	); err != nil {
		return nil
	}

	return nil
}
