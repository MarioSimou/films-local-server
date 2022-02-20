import {
  DynamoDBClient,
  PutItemCommand,
  GetItemCommand,
  AttributeValue,
  DeleteItemCommand,
  ScanCommand,
  ScanCommandInput,
} from '@aws-sdk/client-dynamodb'
import type {Song, Result, GUID, Status} from '@libs/types'

type DynamodbClientFactoryResult = {
  getSongs: ReturnType<typeof getSongs>
  getSong: ReturnType<typeof getSong>
  deleteSong: ReturnType<typeof deleteSong>
  putSong: ReturnType<typeof putSong>
}

const AWS_DYNAMODB_SONGS_TABLE_NAME = process.env.AWS_DYNAMODB_SONGS_TABLE_NAME
const AWS_ENDPOINT = process.env.AWS_ENDPOINT

if (!AWS_DYNAMODB_SONGS_TABLE_NAME) {
  throw new Error('error: dynamodb is not configured correctly')
}

export const ErrDynamoDBItemNotFound = new Error('error: dynamodb item not found')

const bodyToItem = (body: Song): Record<string, AttributeValue> => {
  return Object.entries(body).reduce((item, [attributeName, attributeValue]) => {
    if (!attributeValue) {
      return item
    }
    return {
      ...item,
      [attributeName]: {
        'S': attributeValue,
      },
    }
  }, {})
}

const itemToSong = (item: Record<string, AttributeValue>): Song => {
  return {
    guid: item?.guid?.S,
    description: item?.description?.S,
    title: item?.title?.S,
    createdAt: item?.createdAt?.S,
    image: item?.image?.S,
    imageStatus: item?.imageStatus?.S as Status,
    href: item?.href?.S,
    hrefStatus: item?.hrefStatus?.S as Status,
  }
}

export const putSong =
  (client: DynamoDBClient) =>
  async (body: Song): Promise<Result<Song>> => {
    try {
      const command = new PutItemCommand({
        TableName: AWS_DYNAMODB_SONGS_TABLE_NAME,
        Item: bodyToItem(body),
      })
      await client.send(command)

      const [e, newSong] = await getSong(client)(body.guid)
      if (e) {
        return [e]
      }
      return [undefined, newSong]
    } catch (e) {
      return [e]
    }
  }

export const getSong =
  (client: DynamoDBClient) =>
  async (partitionKey: GUID): Promise<Result<Song>> => {
    try {
      const command = new GetItemCommand({
        TableName: AWS_DYNAMODB_SONGS_TABLE_NAME,
        Key: {
          'guid': {
            'S': partitionKey,
          },
        },
      })
      const {Item: item} = await client.send(command)
      if (!item?.guid) {
        return [ErrDynamoDBItemNotFound]
      }
      return [undefined, itemToSong(item)]
    } catch (e) {
      return [e]
    }
  }

export const getSongs =
  (client: DynamoDBClient) =>
  async (options: Omit<ScanCommandInput, 'TableName'> = {}): Promise<Result<Song[]>> => {
    try {
      const command = new ScanCommand({
        TableName: AWS_DYNAMODB_SONGS_TABLE_NAME,
        ...options,
      })
      const {Items: items} = await client.send(command)
      if (items.length === 0) {
        return [ErrDynamoDBItemNotFound]
      }
      const songs = items.map(itemToSong)
      return [undefined, songs]
    } catch (e) {
      return [e]
    }
  }

export const deleteSong =
  (client: DynamoDBClient) =>
  async (partitionKey: GUID): Promise<[Error | undefined]> => {
    try {
      const command = new DeleteItemCommand({
        TableName: AWS_DYNAMODB_SONGS_TABLE_NAME,
        Key: {
          'guid': {
            'S': partitionKey,
          },
        },
      })
      await client.send(command)
      return [undefined]
    } catch (e) {
      return [e]
    }
  }

const getDynamoDBClient = (): DynamoDBClient => {
  if (!AWS_ENDPOINT) {
    return new DynamoDBClient({})
  }
  return new DynamoDBClient({
    endpoint: AWS_ENDPOINT,
  })
}

export const newDynamoDBClientFactory = (): DynamodbClientFactoryResult => {
  const client = getDynamoDBClient()
  return {
    getSongs: getSongs(client),
    getSong: getSong(client),
    deleteSong: deleteSong(client),
    putSong: putSong(client),
  }
}
