### Setup

```
ENDPOINT_URL=http://localhost:4566
TABLE_NAME=songs-serverless-backend-songs-dev
BUCKET_NAME=songs-serverless-backend-assets-dev
UPLOAD_TOPIC_NAME=songs-serverless-backend-upload-assets-dev
DELETE_TOPIC_NAME=songs-serverless-backend-delete-assets-dev

# DynamoDB
aws dynamodb create-table --table-name $TABLE_NAME --endpoint-url $ENDPOINT_URL --attribute-definitions AttributeName=guid,AttributeType=S --key-schema AttributeName=guid,KeyType=HASH --provisioned-throughput ReadCapacityUnits=3,WriteCapacityUnits=3
aws dynamodb list-tables --endpoint-url $ENDPOINT_URL

# SNS
aws sns create-topic --name $UPLOAD_TOPIC_NAME --endpoint-url $ENDPOINT_URL
aws sns create-topic --name $DELETE_TOPIC_NAME --endpoint-url $ENDPOINT_URL

aws sns list-topics --endpoint-url $ENDPOINT_URL

# S3
aws s3api create-bucket --bucket $BUCKET_NAME --endpoint-url $ENDPOINT_URL --acl public-read
aws s3 ls --endpoint-url $ENDPOINT_URL
```
