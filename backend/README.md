
### Setup 

```
ENDPOINT_URL=http://localhost:4566

# create dynamodb table
aws dynamodb create-table --table-name songs --attribute-definitions AttributeName=GUID,AttributeType=S --key-schema AttributeName=GUID,KeyType=HASH --endpoint-url $ENDPOINT_URL --provisioned-throughput ReadCapacityUnits=3,WriteCapacityUnits=3
```