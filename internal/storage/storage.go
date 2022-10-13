package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
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
	listContents("/github/workspace")
	listContents("/github/home")
	listContents("/github/workflow")

	filePath := s.CodePath
	if _, ok := os.LookupEnv("GITHUB_SHA"); ok {
		filePath = fmt.Sprintf("/github/workspace/%s", filePath)
	}

	log.Printf("trying to open file %s\n", filePath)

	dat, err := os.ReadFile(filePath)
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

func listContents(path string) {
	log.Printf("listing contents for path %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
