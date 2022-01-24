package awsUtils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type PublishUploadMessage struct {
	Key   string          `json:"key"`
	Body  []byte          `json:"body"`
	Field repoTypes.Field `json:"field"`
}

func NewPublishMessageUpload(subject, key string, field repoTypes.Field, body []byte) (*sns.PublishInput, error) {
	var payload = PublishUploadMessage{
		Key:   key,
		Body:  []byte(base64.StdEncoding.EncodeToString(body)),
		Field: field,
	}
	var payloadBf []byte
	var e error
	if payloadBf, e = json.Marshal(&payload); e != nil {
		return nil, e
	}

	return &sns.PublishInput{
		TopicArn: aws.String(os.Getenv("AWS_SNS_UPLOAD_TOPIC_ARN")),
		Subject:  &subject,
		Message:  aws.String(string(payloadBf)),
	}, nil
}

func UploadOne(ctx context.Context, client *sns.Client, key string, field repoTypes.Field, body []byte) (*string, error) {
	var e error
	var publishInput *sns.PublishInput
	var publishOutput *sns.PublishOutput

	if publishInput, e = NewPublishMessageUpload("upload", key, field, body); e != nil {
		return nil, e
	}

	if publishOutput, e = client.Publish(ctx, publishInput); e != nil {
		return nil, e
	}
	return publishOutput.MessageId, nil
}
