import { formatJSONResponse, newHandler, ErrNotFound, StatusNotFound, StatusBadRequest,StatusInternalServerError, StatusOk } from '@libs/api-gateway'
import type { Handler, GUIDPathParameters, Song } from '@libs/types'
import * as Joi from 'joi'
import { findOneSong, putOneSong } from '@libs/dynamodb'

const pathParametersSchema = Joi.object({
    guid: Joi.string().required().uuid(),
})

const bodySchema = Joi.object({
    name: Joi.string(),
    description: Joi.string(),
})

const baseHandler: Handler<Partial<Pick<Song, 'description' | 'name'>>, GUIDPathParameters> = async (event) => {
    const {error: pathValidationError} = pathParametersSchema.validate(event.pathParameters)
    if(pathValidationError){
        return formatJSONResponse(StatusBadRequest, pathValidationError.message)
    }

    const {error: bodyValidationError} = bodySchema.validate(event.body)
    if(bodyValidationError){
        return formatJSONResponse(StatusBadRequest, bodyValidationError.message)
    }

    const [findOneSongError, currentSong] = await findOneSong(event.pathParameters.guid)
    if(findOneSongError){
        if(findOneSongError === ErrNotFound){
            return formatJSONResponse(StatusNotFound, findOneSongError.message)
        }
        return formatJSONResponse(StatusInternalServerError, findOneSongError.message)
    }
    
    const [putOneSongError, newSong] = await putOneSong(currentSong.guid, {...currentSong,...event.body})
    if(putOneSongError){
        return formatJSONResponse(StatusInternalServerError, putOneSongError.message)
    }

    return formatJSONResponse(StatusOk, newSong)
}

const handler = newHandler(baseHandler)

export default handler