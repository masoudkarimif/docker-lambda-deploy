package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	n := &notification.Notification{
		Hook: os.Getenv("INPUT_SLACK_HOOK"),
	}

	if err := n.SendInProgressMsg(ctx); err != nil {
		log.Println("sending slack failed, %s", err.Error())
	}

	cfg, err := getAwsConfig(ctx)
	if err != nil {
		log.Println("get config failed, %s", err.Error())
		os.Exit(1)
	}

	s := &storage.Storage{
		BucketName: os.Getenv("INPUT_BUCKET"),
		Key:        os.Getenv("INPUT_KEY"),
		CodePath:   os.Getenv("INPUT_CODE_PATH"),
	}
	s.Initialize(cfg)

	f := &function.Function{
		BucketName: os.Getenv("INPUT_BUCKET"),
		Key:        os.Getenv("INPUT_KEY"),
		Name:       os.Getenv("INPUT_FUNCTION_NAME"),
	}
	f.Initialize(cfg)

	a := &action.Action{}
	if err := a.Run(ctx, f, s); err != nil {
		fmt.Fprintf(os.Stdout, "run exited, %s\n", err)
		if err := n.SendFailedMsg(ctx); err != nil {
			log.Println("sending slack failed, %s", err.Error())
		}
		os.Exit(1)
	}

	if err := n.SendSucceededMsg(ctx); err != nil {
		log.Println("sending slack failed, %s", err.Error())
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
	return len(os.Getenv("INPUT_ACCESS_KEY")) > 0 && len(os.Getenv("INPUT_SECRET_KEY")) > 0
}