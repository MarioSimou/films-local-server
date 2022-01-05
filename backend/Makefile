functionName ?= GetSongs
env ?= dev
name = songs-backend
stackBucketName = songs-local-server-bucket

validate:
	@cfn-lint template.yaml

create-stack-bucket:
	@aws s3api create-bucket --bucket $(stackBucketName)

build-functions:
	@go build -o ./.aws-sam/build/PostSong/main ./functions/postSong/main.go
	@go build -o ./.aws-sam/build/GetSongs/main ./functions/getSongs/main.go
	@go build -o ./.aws-sam/build/GetSong/main ./functions/getSong/main.go
	@go build -o ./.aws-sam/build/DeleteSong/main ./functions/deleteSong/main.go
	@go build -o ./.aws-sam/build/PutSong/main ./functions/putSong/main.go
	@go build -o ./.aws-sam/build/UploadSong/main ./functions/uploadSong/main.go
	@go build -o ./.aws-sam/build/PutSongImage/main ./functions/putSongImage/main.go
	@go build -o ./.aws-sam/build/DeleteSongImage/main ./functions/deleteSongImage/main.go

watch-build-functions:
	@reflex -s -r '\.go$$' -- make build-functions

test:
	@sam build --parallel
	@sam local invoke $(functionName) --parameter-overrides Environment='prod'

test-api:
	@rm -rf .aws-sam
	@sam build --parallel
	@sam local start-api --parameter-overrides Environment=$(env) AWSS3Endpoint='http://s3:3003' AWSDynamoDBEndpoint='http://dynamodb:3001' --docker-network songs-local-server-backend-network

packages:
	@sam package --s3-bucket $(stackBucketName)

deploy:
	@sam build --parallel
	@sam deploy --stack-name $(name) --s3-bucket $(stackBucketName) --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM --parameter-overrides Environment=prod

cleanup:
	@sam delete --stack-name $(name) 