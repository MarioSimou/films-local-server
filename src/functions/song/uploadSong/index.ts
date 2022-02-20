import type {Handler, GUID} from '@libs/types'
import {
  newHandler,
  newResponse,
  StatusInternalServerError,
  StatusBadRequest,
  StatusOk,
  StatusNotFound,
  StatusUnsupportedMediaType,
  ErrUnsupportedMediaType,
  parseMultipartFormData,
} from '@libs/api-gateway'
import aws, {ErrDynamoDBItemNotFound} from '@libs/aws'
import Joi from 'joi'
import {extname} from 'path'

type PathParamsSchema = {guid: GUID}
const pathParamsSchema = Joi.object<PathParamsSchema>({
  guid: Joi.string().required().uuid(),
})

const baseHandler: Handler<string, PathParamsSchema> = async event => {
  const {error: pathParametersValidationError} = pathParamsSchema.validate(event.pathParameters)
  if (pathParametersValidationError) {
    return newResponse(StatusBadRequest, pathParametersValidationError)
  }
  const [parseMultipartFormDataError, fileData] = parseMultipartFormData(event)
  if (parseMultipartFormDataError) {
    if (parseMultipartFormDataError === ErrUnsupportedMediaType) {
      return newResponse(StatusUnsupportedMediaType, parseMultipartFormDataError)
    }
    return newResponse(StatusBadRequest, parseMultipartFormDataError)
  }
  const {guid} = event.pathParameters

  const [getSongError, currentSong] = await aws.dynamoDB.getSong(guid)
  if (getSongError) {
    if (getSongError === ErrDynamoDBItemNotFound) {
      return newResponse(StatusNotFound, getSongError)
    }
    return newResponse(StatusInternalServerError, getSongError)
  }

  const field = 'href'
  const [publishError] = await aws.sns.publishUploadAssetMessage({
    song: currentSong,
    key: `${currentSong.guid}/${field}${extname(fileData.filename)}`,
    asset: fileData.content.toString('base64'),
    contentType: fileData.contentType,
    field,
  })
  if (publishError) {
    return newResponse(StatusInternalServerError, publishError)
  }
  const [putSongError, newSong] = await aws.dynamoDB.putSong({
    ...currentSong,
    hrefStatus: 'pending',
  })
  if (putSongError) {
    return newResponse(StatusInternalServerError, putSongError)
  }

  return newResponse(StatusOk, newSong)
}

export const handler = newHandler(baseHandler)
