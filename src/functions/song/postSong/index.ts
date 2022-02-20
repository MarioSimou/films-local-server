import type {Handler, Song} from '@libs/types'
import {
  newHandler,
  newResponse,
  StatusOk,
  StatusBadRequest,
  StatusCreated,
  StatusInternalServerError,
} from '@libs/api-gateway'
import Joi from 'joi'
import {v4 as uuidv4} from 'uuid'
import aws from '@libs/aws'

type BodySchema = Pick<Song, 'title' | 'description'>
const bodySchema = Joi.object<BodySchema>({
  title: Joi.string().required(),
  description: Joi.string().max(1000),
})

const baseHandler: Handler<BodySchema> = async event => {
  const {error: validationError} = await bodySchema.validate(event.body)

  if (validationError) {
    return newResponse(StatusBadRequest, validationError)
  }

  const guid = uuidv4()

  const [putSongError, song] = await aws.dynamoDB.putSong({
    guid,
    createdAt: new Date().toISOString(),
    ...event.body,
  })
  if (putSongError) {
    return newResponse(StatusInternalServerError, putSongError)
  }

  return newResponse(StatusCreated, song, {
    'Location': `/api/v1/songs/${song.guid}`,
  })
}

export const handler = newHandler(baseHandler)
