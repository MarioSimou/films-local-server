package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/MarioSimou/songs-local-server/backend/packages/awsUtils"
	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	awsClients *awsUtils.AWSClients
)

func init() {
	var e error
	var ctx = context.Background()
	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("error: %v\n", e)
	}
}

type UploadMessage struct {
	Key   string          `json:"key"`
	Body  string          `json:"body"`
	Field repoTypes.Field `json:"field"`
}

func getGUIDFromKey(key string) string {
	return strings.Split(key, "/")[0]
}

func Handler(ctx context.Context, event events.SNSEvent) error {
	var record = event.Records[0]
	var message UploadMessage
	var body []byte

	if e := json.NewDecoder(strings.NewReader(record.SNS.Message)).Decode(&message); e != nil {
		return e
	}

	var guid = getGUIDFromKey(message.Key)
	var currentSong *repoTypes.Song
	var e error
	if currentSong, e = awsUtils.GetOneSong(ctx, guid, awsClients.DynamoDB); e != nil {
		return e
	}

	if body, e = base64.StdEncoding.DecodeString(message.Body); e != nil {
		return e
	}

	var uploadObjectOutput *manager.UploadOutput
	if uploadObjectOutput, e = awsUtils.UploadObject(ctx, message.Key, body, types.StorageClassStandard, awsClients.S3Uploader); e != nil {
		return e
	}

	switch message.Field {
	case repoTypes.LocationField:
		currentSong.Location = uploadObjectOutput.Location
		currentSong.LocationStatus = repoTypes.Done
	case repoTypes.ImageField:
		currentSong.Image = uploadObjectOutput.Location
		currentSong.ImageStatus = repoTypes.Done
	default:
		return fmt.Errorf("error: field does not match")
	}

	if _, e := awsUtils.PutSong(ctx, *currentSong, awsClients.DynamoDB); e != nil {
		return e
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
