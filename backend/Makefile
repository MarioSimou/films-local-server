functionName ?= GetFilms
env ?= stage
name = films-local-server-backend
bucketName = films-local-server-backend-bucket

validate:
	@cfn-lint template.yaml

create-bucket:
	@aws s3api create-bucket --bucket $(bucketName) 

test:
	@sam build --parallel
	@sam local invoke $(functionName) --parameter-overrides Environment=$(env)

test-api:
	@sam build --parallel
	@sam local start-api --parameter-overrides Environment=$(env)

packages:
	@sam package --s3-bucket $(bucketName)

deploy:
	@make validate
	@sam build --parallel
	@sam deploy --stack-name $(name) --s3-bucket $(bucketName) --capabilities CAPABILITY_IAM --parameter-overrides Environment=prod

cleanup:
	@sam delete --stack-name $(name) 