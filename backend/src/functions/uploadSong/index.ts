import { formatJSONResponse, newHandler, ErrNotFound, StatusNotFound, StatusBadRequest,StatusInternalServerError, StatusOk, isMultipart, ErrNoMultipart, parseMultipart } from '@libs/api-gateway'
import type { Handler, GUIDPathParameters, Song } from '@libs/types'
import * as Joi from 'joi'
import { findOneSong, putOneSong } from '@libs/dynamodb'
import { uploadOne } from '@libs/s3'
import { extname } from 'path'

const pathParametersSchema = Joi.object({
    guid: Joi.string().required().uuid(),
})

const baseHandler: Handler<any, GUIDPathParameters> = async (event) => {
    const {error: validationError} = pathParametersSchema.validate(event.pathParameters)

    if(validationError){
        return formatJSONResponse(StatusBadRequest, validationError.message)
    }

    if(!isMultipart(event.headers)){
        return formatJSONResponse(StatusBadRequest, ErrNoMultipart.message)
    }

    const [parseMultipartError, multipart] = parseMultipart(event)
    if(parseMultipartError){
        return formatJSONResponse(StatusBadRequest, parseMultipartError.message)
    }

    const [findOneSongError, currentSong] = await findOneSong(event.pathParameters.guid)
    if(findOneSongError){
        if(findOneSongError === ErrNotFound){
            return formatJSONResponse(StatusNotFound, findOneSongError.message)
        }
        return formatJSONResponse(StatusInternalServerError, findOneSongError.message)
    }

    const key = `${currentSong.guid}/song${extname(multipart.filename)}` 
    const [uploadOneError, songLocation] = await uploadOne(key, multipart.content, {ContentType: multipart.contentType})
    if(uploadOneError){
        return formatJSONResponse(StatusInternalServerError, uploadOneError.message)
    }
    
    const [putOneSongError, newSong] = await putOneSong(currentSong.guid, {...currentSong, href: songLocation})
    if(putOneSongError){
        return formatJSONResponse(StatusInternalServerError, putOneSongError.message)
    }

    return formatJSONResponse(StatusOk, newSong)
}

const handler = newHandler(baseHandler)

export default handler