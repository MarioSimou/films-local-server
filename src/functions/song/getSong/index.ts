import type {Handler, GUID} from '@libs/types'
import {
  newHandler,
  newResponse,
  StatusOk,
  StatusNotFound,
  StatusInternalServerError,
  StatusBadRequest,
} from '@libs/api-gateway'
import aws, {ErrDynamoDBItemNotFound} from '@libs/aws'
import Joi from 'joi'

type PathParamsSchema = {guid: GUID}
const pathParamsSchema = Joi.object<PathParamsSchema>({
  guid: Joi.string().required().uuid(),
})

const baseHandler: Handler<unknown, PathParamsSchema> = async event => {
  const {error: validationError} = pathParamsSchema.validate(event.pathParameters)
  if (validationError) {
    return newResponse(StatusBadRequest, validationError)
  }

  const {guid} = event.pathParameters

  const [getSongError, song] = await aws.dynamoDB.getSong(guid)
  if (getSongError) {
    if (getSongError === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, getSongError)
    }
    return newResponse(StatusInternalServerError, getSongError)
  }

  return newResponse(StatusOk, song)
}

export const handler = newHandler(baseHandler)
