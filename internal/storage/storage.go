package storage

import (
	"bytes"
	"context"
	"fmt"
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
	dat, err := s.ReadFile(s.CodePath)
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

	log.Printf("successfully wrote zip file to s3://%s/%s\n", s.BucketName, s.Key)

	return nil
}

func (*Storage) ReadFile(filePath string) ([]byte, error) {
	if _, ok := os.LookupEnv("GITHUB_SHA"); ok {
		filePath = fmt.Sprintf("/github/workspace/%s", filePath)
	}

	log.Printf("trying to open file %s\n", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
