import {SNSClient, PublishCommand, PublishResponse} from '@aws-sdk/client-sns'
import type {Result, Song} from '@libs/types'

const AWS_ENDPOINT = process.env.AWS_ENDPOINT
const AWS_SNS_UPLOAD_TOPIC_ARN = process.env.AWS_SNS_UPLOAD_TOPIC_ARN
const AWS_SNS_DELETE_TOPIC_ARN = process.env.AWS_SNS_DELETE_TOPIC_ARN

if (!AWS_SNS_UPLOAD_TOPIC_ARN || !AWS_SNS_DELETE_TOPIC_ARN) {
  throw new Error('error: sns configuration is invalid')
}

export const TOPICS = {
  UPLOAD: AWS_SNS_UPLOAD_TOPIC_ARN,
  DELETE: AWS_SNS_DELETE_TOPIC_ARN,
}

const getSNSClient = (): SNSClient => {
  if (!AWS_ENDPOINT) {
    return new SNSClient({})
  }
  return new SNSClient({endpoint: AWS_ENDPOINT})
}

const publish =
  (client: SNSClient) =>
  async (subject: string, message: string, topicArn: string): Promise<Result<PublishResponse>> => {
    console.log('topics: ', AWS_SNS_UPLOAD_TOPIC_ARN, AWS_SNS_DELETE_TOPIC_ARN)
    try {
      const command = new PublishCommand({
        TopicArn: topicArn,
        Subject: subject,
        Message: message,
      })
      const result = await client.send(command)
      return [undefined, result]
    } catch (e) {
      return [e]
    }
  }

type Field = 'href' | 'image'
export type PublishUploadAssetMessage = {song: Song; asset: string; key: string; field: Field; contentType: string}
const publishUploadAssetMessage = (client: SNSClient) => async (message: PublishUploadAssetMessage) => {
  return publish(client)('upload', JSON.stringify(message), TOPICS.UPLOAD)
}

export type PublishDeleteAssetMessage = {song: Song; field: Field}
const publishDeleteAssetMessage = (client: SNSClient) => async (message: PublishDeleteAssetMessage) => {
  return publish(client)('delete', JSON.stringify(message), TOPICS.DELETE)
}

type SNSClientFactoryResult = {
  publish: ReturnType<typeof publish>
  publishUploadAssetMessage: ReturnType<typeof publishUploadAssetMessage>
  publishDeleteAssetMessage: ReturnType<typeof publishDeleteAssetMessage>
}
export const newSNSClientFactory = (): SNSClientFactoryResult => {
  const client = getSNSClient()
  return {
    publish: publish(client),
    publishDeleteAssetMessage: publishDeleteAssetMessage(client),
    publishUploadAssetMessage: publishUploadAssetMessage(client),
  }
}
