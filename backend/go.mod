module github.com/others/songs-local-server-sls

go 1.15

require (
	github.com/aws/aws-lambda-go v1.6.0
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.13.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.6.0
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.9.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.13.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.24.0
	github.com/envoyproxy/protoc-gen-validate v0.6.3
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/jinzhu/copier v0.3.5 // indirect
	go.mongodb.org/mongo-driver v1.8.2 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)
