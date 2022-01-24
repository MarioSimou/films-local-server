package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MarioSimou/songs-local-server/backend/packages/awsUtils"
	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
	"github.com/MarioSimou/songs-local-server/backend/packages/utils"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/copier"
)

var (
	awsClients *awsUtils.AWSClients
)

type bodyBinding struct {
	Name        string    `json:"name" validate:"required_without_all=Description"`
	Description string    `json:"description" validate:"max=1000"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func init() {
	var e error
	var ctx = context.Background()

	if awsClients, e = awsUtils.NewAWSClients(ctx); e != nil {
		log.Fatalf("Error: %v\n", e)

	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var songGUID = req.PathParameters["guid"]
	var validate = utils.NewValidator()
	var body bodyBinding
	var newSong *repoTypes.Song
	var currentSong *repoTypes.Song
	var e error

	if e := utils.DecodeEventBody(req.Body, &body); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	if e := validate.Var(songGUID, "required,uuid4"); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}

	if e := validate.Struct(body); e != nil {
		return utils.NewAPIResponse(http.StatusBadRequest, e), nil
	}
	body.UpdatedAt = time.Now()

	if currentSong, e = awsUtils.GetOneSong(ctx, songGUID, awsClients.DynamoDB); e != nil {
		if e == repoTypes.ErrSongNotFound {
			return utils.NewAPIResponse(http.StatusNotFound, e), nil
		}
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if e := copier.CopyWithOption(currentSong, &body, copier.Option{IgnoreEmpty: true}); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	if newSong, e = awsUtils.PutSong(ctx, *currentSong, awsClients.DynamoDB); e != nil {
		return utils.NewAPIResponse(http.StatusInternalServerError, e), nil
	}

	return utils.NewAPIResponse(http.StatusOK, newSong), nil
}

func main() {
	runtime.Start(handler)
}
