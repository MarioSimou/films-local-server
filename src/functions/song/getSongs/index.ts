import type {Handler} from '@libs/types'
import {newHandler, newResponse, StatusOk, StatusNotFound, StatusInternalServerError} from '@libs/api-gateway'
import aws, {ErrDynamoDBItemNotFound} from '@libs/aws'

const baseHandler: Handler<unknown, unknown, {limit?: number}> = async event => {
  const {limit} = event.queryStringParameters ?? {limit: 20}

  const [getSongsError, songs] = await aws.dynamoDB.getSongs({Limit: limit})
  if (getSongsError) {
    if (getSongsError === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, getSongsError)
    }
    return newResponse(StatusInternalServerError, getSongsError)
  }

  return newResponse(StatusOk, songs)
}

export const handler = newHandler(baseHandler)
