import {newDynamoDBClientFactory} from './dynamodb'
import {newS3ClientFactory} from './s3'
import {newSNSClientFactory} from './sns'

export * from './dynamodb'
export * from './s3'
export * from './sns'

export type AWSClients = {
  dynamoDB: ReturnType<typeof newDynamoDBClientFactory>
  s3: ReturnType<typeof newS3ClientFactory>
  sns: ReturnType<typeof newSNSClientFactory>
}

export default {
  dynamoDB: newDynamoDBClientFactory(),
  s3: newS3ClientFactory(),
  sns: newSNSClientFactory(),
} as AWSClients
