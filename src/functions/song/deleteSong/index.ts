import type {Handler, GUID} from '@libs/types'
import {
  newHandler,
  newResponse,
  StatusNotFound,
  StatusInternalServerError,
  StatusBadRequest,
  StatusNoContent,
} from '@libs/api-gateway'
import aws, {ErrDynamoDBItemNotFound, TOPICS} from '@libs/aws'
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

  const [e, currentSong] = await aws.dynamoDB.getSong(guid)
  if (e) {
    if (e === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, e)
    }
    return newResponse(StatusInternalServerError, e)
  }

  if (currentSong.href) {
    await aws.sns.publishDeleteAssetMessage({song: currentSong, field: 'href'})
  }
  if (currentSong.image) {
    await aws.sns.publishDeleteAssetMessage({song: currentSong, field: 'image'})
  }

  const [deleteSongError] = await aws.dynamoDB.deleteSong(guid)
  if (deleteSongError) {
    if (deleteSongError === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, deleteSongError)
    }
    return newResponse(StatusInternalServerError, deleteSongError)
  }

  return newResponse(StatusNoContent)
}

export const handler = newHandler(baseHandler)
