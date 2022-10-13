package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

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
	fullPath := getFileFullPath(filePath)
	log.Printf("trying to open file %s\n", fullPath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getFileFullPath(filePath string) string {
	var filePrefix string
	if _, ok := os.LookupEnv("GITHUB_SHA"); ok {
		filePrefix = "/github/workspace/"
	}

	workingDirectory := os.Getenv("INPUT_WORKING_DIRECTORY")
	listFilesInDirectory(path.Join(filePrefix, workingDirectory))
	fullPath := path.Join(filePrefix, workingDirectory, filePath)

	return fullPath
}

func listFilesInDirectory(directoryPath string) {
	log.Printf("Listing files in directory %s\n", directoryPath)
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
