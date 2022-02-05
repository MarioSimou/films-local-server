## Setup

```
ENDPOINT_URL=http://localhost:4566
TABLE_NAME=songs-local-server-backend-songs-dev
BUCKET_NAME=songs-local-server-backend-asssets-bucket-dev

aws dynamodb create-table --endpoint-url $ENDPOINT_URL --table-name $TABLE_NAME --key-schema AttributeName=guid,KeyType=HASH  --attribute-definitions AttributeName=guid,AttributeType=S --provisioned-throughput ReadCapacityUnits=3,WriteCapacityUnits=3

aws s3api create-bucket --bucket $BUCKET_NAME  --acl public-read --endpoint-url $ENDPOINT_URL
```
