import type {Handler, GUID} from '@libs/types'
import {
  newHandler,
  newResponse,
  StatusInternalServerError,
  StatusBadRequest,
  StatusOk,
  StatusNotFound,
} from '@libs/api-gateway'
import aws, {ErrDynamoDBItemNotFound} from '@libs/aws'
import Joi from 'joi'

type BodySchema = {
  title?: string
  description?: string
}
const bodySchema = Joi.object<BodySchema>({
  title: Joi.string().when('description', {
    is: Joi.any().empty(),
    then: Joi.string().required(),
  }),
  description: Joi.string(),
})

type PathParamsSchema = {guid: GUID}
const pathParamsSchema = Joi.object<PathParamsSchema>({
  guid: Joi.string().required().uuid(),
})

const baseHandler: Handler<BodySchema, PathParamsSchema> = async event => {
  const {error: pathParametersValidationError} = pathParamsSchema.validate(event.pathParameters)
  if (pathParametersValidationError) {
    return newResponse(StatusBadRequest, pathParametersValidationError)
  }

  const {error: bodyValidationError} = bodySchema.validate(event.body)
  if (bodyValidationError) {
    return newResponse(StatusBadRequest, bodyValidationError)
  }

  const {guid} = event.pathParameters

  const [getSongError, currentSong] = await aws.dynamoDB.getSong(guid)
  if (getSongError) {
    if (getSongError === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, getSongError)
    }
    return newResponse(StatusInternalServerError, getSongError)
  }

  const [putSongError, newSong] = await aws.dynamoDB.putSong({
    ...currentSong,
    title: event.body.title ?? currentSong.title,
    description: event.body.description ?? currentSong.description,
  })

  if (putSongError) {
    return newResponse(StatusInternalServerError, putSongError)
  }

  return newResponse(StatusOk, newSong)
}

export const handler = newHandler(baseHandler)
