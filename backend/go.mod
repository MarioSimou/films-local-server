module github.com/MarioSimou/songs-local-server/backend

go 1.15

require (
	github.com/aws/aws-lambda-go v1.27.1
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.11.1
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.4.5
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.7.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.11.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.22.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.15.0
	github.com/aws/smithy-go v1.10.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/google/uuid v1.3.0
	github.com/jinzhu/copier v0.3.4
)
