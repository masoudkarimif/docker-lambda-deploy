package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/masoudkarimif/docker-lambda-deploy/internal/action"
	"github.com/masoudkarimif/docker-lambda-deploy/internal/function"
	"github.com/masoudkarimif/docker-lambda-deploy/internal/notification"
	"github.com/masoudkarimif/docker-lambda-deploy/internal/storage"
)

type cred struct{}

func (cred) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     os.Getenv("INPUT_AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("INPUT_AWS_SECRET_ACCESS_KEY"),
	}, nil
}

func main() {
	ctx := context.Background()

	cfg, err := getAwsConfig(ctx)
	if err != nil {
		log.Printf("get config failed, %s", err.Error())
		os.Exit(1)
	}

	s := &storage.Storage{
		BucketName: os.Getenv("INPUT_S3_BUCKET"),
		Key:        os.Getenv("INPUT_S3_KEY"),
		CodePath:   os.Getenv("INPUT_CODE_PATH"),
	}
	s.Initialize(cfg)

	// resolve function name
	if fnName, ok := os.LookupEnv("INPUT_FUNCTION_NAME"); !ok || len(fnName) == 0 {
		fnNameByte, err := s.ReadFile(".function_name")
		if err != nil {
			log.Fatal(err.Error())
		}

		fnName = strings.TrimSpace(string(fnNameByte))
		os.Setenv("INPUT_FUNCTION_NAME", fnName)
	}

	n := &notification.Notification{
		Hook: os.Getenv("INPUT_SLACK_HOOK"),
	}

	if err := n.SendInProgressMsg(ctx); err != nil {
		log.Printf("sending slack failed, %s", err.Error())
	}

	f := &function.Function{
		BucketName: os.Getenv("INPUT_S3_BUCKET"),
		Key:        os.Getenv("INPUT_S3_KEY"),
		Name:       os.Getenv("INPUT_FUNCTION_NAME"),
	}
	f.Initialize(cfg)

	a := &action.Action{}
	if err := a.Run(ctx, f, s); err != nil {
		fmt.Fprintf(os.Stdout, "run exited, %s\n", err)
		if err := n.SendFailedMsg(ctx); err != nil {
			log.Printf("sending slack msg failed, %s", err.Error())
		}
		os.Exit(1)
	}

	if err := n.SendSucceededMsg(ctx); err != nil {
		log.Printf("sending slack msg failed, %s", err.Error())
	}
}

func getAwsConfig(ctx context.Context) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if keysPresent() {
		c := cred{}
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("INPUT_AWS_REGION")),
			config.WithCredentialsProvider(c))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("INPUT_AWS_REGION")))
	}

	return cfg, err
}

func keysPresent() bool {
	return len(os.Getenv("INPUT_AWS_ACCESS_KEY_ID")) > 0 && len(os.Getenv("INPUT_AWS_SECRET_ACCESS_KEY")) > 0
}
