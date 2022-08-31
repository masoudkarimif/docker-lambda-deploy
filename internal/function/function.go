package function

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Function struct {
	Name       string
	BucketName string
	Key        string
	client     *lambda.Client
}

func (f *Function) Initialize(cfg aws.Config) {
	f.client = lambda.NewFromConfig(cfg)
}

func (f *Function) UpdateCode(ctx context.Context) error {
	if _, err := f.client.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(f.Name),
		S3Bucket:     aws.String(f.BucketName),
		S3Key:        aws.String(f.Key),
	}); err != nil {
		return err
	}

	log.Printf("successfully updated lambda function %s\n", f.Name)

	return nil
}
