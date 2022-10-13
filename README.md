# docker-lambda-deploy

![build](https://github.com/masoudkarimif/docker-lambda-deploy/actions/workflows/build.yml/badge.svg)
![golangci-lint](https://github.com/masoudkarimif/docker-lambda-deploy/actions/workflows/golangci-lint.yml/badge.svg)

Update Lambda functions + Slack notifications.


## Run
```bash
docker run --rm \                                                                1 ↵ mkf@MKFs-MacBook-Pro
  -e INPUT_AWS_ACCESS_KEY_ID=xxxxxxx \
  -e INPUT_AWS_SECRET_ACCESS_KEY=xxxxxx \
  -e INPUT_AWS_REGION=us-east-1 \
  -e INPUT_FUNCTION_NAME=my-function \
  -e INPUT_S3_KEY=artifact.zip \
  -e INPUT_S3_BUCKET=devops-bucket \
  -e INPUT_CODE_PATH=./build/artifact.zip \
  -e INPUT_SLACK_HOOK=https://xxxxxxxxxxxx \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  masoudkf/docker-lambda-deploy
```

## AWS permissions needed
`s3:PutObject` and `lambda:UpdateFunctionCode`

