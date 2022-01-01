functionName ?= GetSongs
env ?= stage
name = songs-local-server-backend
bucketName = songs-local-server-backend-bucket

validate:
	@cfn-lint template.yaml

create-bucket:
	@aws s3api create-bucket --bucket $(bucketName)

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
	@sam local invoke $(functionName) --parameter-overrides Environment=$(env)

test-api:
	@rm -rf .aws-sam
	@sam build --parallel
	@sam local start-api --parameter-overrides Environment=$(env) --docker-network songs-local-server-backend-network

packages:
	@sam package --s3-bucket $(bucketName)

deploy:
	@make validate
	@sam build --parallel
	@sam deploy --stack-name $(name) --s3-bucket $(bucketName) --capabilities CAPABILITY_IAM --parameter-overrides Environment=prod

cleanup:
	@sam delete --stack-name $(name) 