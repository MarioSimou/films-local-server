import { formatJSONResponse, StatusBadRequest, StatusOk,StatusNotFound, StatusInternalServerError, newHandler, ErrNotFound } from '@libs/api-gateway'
import type { Handler, GUIDPathParameters, Song} from '@libs/types'
import { findOneSong } from '@libs/dynamodb'
import * as Joi from 'joi'

const schema = Joi.object({
    guid: Joi.string().required().uuid()
})

const baseHandler: Handler<unknown, GUIDPathParameters> = async (event) => {
    const {error: validationError} = schema.validate(event.pathParameters)
    if(validationError){
        return formatJSONResponse(StatusBadRequest, validationError.message)
    }

    const [e, song] = await findOneSong(event.pathParameters.guid)
    if(e){
        if(e === ErrNotFound){
            return formatJSONResponse(StatusNotFound, e.message)
        }
        return formatJSONResponse(StatusInternalServerError, e.message)
    }

    return formatJSONResponse(StatusOk, song)
}

const handler = newHandler(baseHandler)
export default handler