functionName ?= GetFilms
env ?= stage
name = films-local-server-backend
bucketName = films-local-server-backend-bucket

validate:
	@cfn-lint template.yaml

create-bucket:
	@aws s3api create-bucket --bucket $(bucketName) 

build-functions:
	@go build -o ./.aws-sam/build/PostFilm/main ./functions/postFilm/main.go
	@go build -o ./.aws-sam/build/GetFilms/main ./functions/getFilms/main.go
	@go build -o ./.aws-sam/build/GetFilm/main ./functions/getFilm/main.go
	@go build -o ./.aws-sam/build/DeleteFilm/main ./functions/deleteFilm/main.go
	@go build -o ./.aws-sam/build/PutFilm/main ./functions/putFilm/main.go
	@go build -o ./.aws-sam/build/UploadFilm/main ./functions/uploadFilm/main.go
	@go build -o ./.aws-sam/build/PutFilmSubtitle/main ./functions/putFilmSubtitle/main.go
	@go build -o ./.aws-sam/build/DeleteFilmSubtitle/main ./functions/deleteFilmSubtitle/main.go

watch-build-functions:
	@reflex -s -r '\.go$$' -- make build-functions

test:
	@sam build --parallel
	@sam local invoke $(functionName) --parameter-overrides Environment=$(env)

test-api:
	@rm -rf .aws-sam
	@sam build --parallel
	@sam local start-api --parameter-overrides Environment=$(env) --docker-network films-local-server-backend-network

packages:
	@sam package --s3-bucket $(bucketName)

deploy:
	@make validate
	@sam build --parallel
	@sam deploy --stack-name $(name) --s3-bucket $(bucketName) --capabilities CAPABILITY_IAM --parameter-overrides Environment=prod

cleanup:
	@sam delete --stack-name $(name) 