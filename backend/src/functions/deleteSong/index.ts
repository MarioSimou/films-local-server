import { formatJSONResponse, newHandler,StatusNoContent, StatusBadRequest, StatusNotFound,StatusInternalServerError, ErrNotFound } from '@libs/api-gateway'
import type { Handler, GUIDPathParameters } from '@libs/types'
import { findOneSong, deleteOne } from '@libs/dynamodb'
import { deleteOne as deleteOneObject } from '@libs/s3'
import * as Joi from 'joi'
import { extname } from 'path'

const schema = Joi.object({
    guid: Joi.string().required().uuid()
})

const baseHandler: Handler<unknown, GUIDPathParameters> = async (event) => {
    const {error: validationError} = schema.validate(event.pathParameters)

    if(validationError){
        return formatJSONResponse(StatusBadRequest, validationError.message)
    }

    const [findOneSongError, song] = await findOneSong(event.pathParameters.guid)
    if(findOneSongError){
        if(findOneSongError === ErrNotFound){
            return formatJSONResponse(StatusNotFound, findOneSongError.message)
        }
        return formatJSONResponse(StatusInternalServerError, findOneSongError.message)
    }

    if(song.href){
        const hrefKey = `${song.guid}/song${extname(song.href)}`
        const [deleteHrefError] = await deleteOneObject(hrefKey)
        if(deleteHrefError){
            return formatJSONResponse(StatusInternalServerError, deleteHrefError.message)
        }        
    }
    if(song.image){
        const imageKey = `${song.guid}/image${extname(song.image)}`
        const [deleteImageError] = await deleteOneObject(imageKey)
        if(deleteImageError){
            return formatJSONResponse(StatusInternalServerError, deleteImageError.message)
        }        
    }

    const [deleteOneError] = await deleteOne(song.guid)
    if(deleteOneError){
        return formatJSONResponse(StatusInternalServerError, findOneSongError.message)
    }

    return formatJSONResponse(StatusNoContent)
}

const handler = newHandler(baseHandler)

export default handler