import { formatJSONResponse, StatusOk, StatusBadRequest, newHandler, StatusInternalServerError } from '@libs/api-gateway'
import type { APIGatewayProxyRequest, APIGatewayProxyResponse, Song  } from '@libs/types'
import * as Joi from 'joi'
import { postSong } from '@libs/dynamodb'


const schema: Joi.ObjectSchema<Pick<Song, 'name' | 'description'>> = Joi.object({
    name: Joi.string().required(),
    description: Joi.string(),
})

const baseHandler = async (event: APIGatewayProxyRequest<{name: string, description?: string}>): Promise<APIGatewayProxyResponse> => {
    const { error: validationError } = schema.validate(event.body)
    if(validationError){
        return formatJSONResponse(StatusBadRequest, validationError.message)
    }

    const [e, song] = await postSong(event.body)
    if(e){
        return formatJSONResponse(StatusInternalServerError, e.message)
    }

    return formatJSONResponse(StatusOk, song)
}

const handler = newHandler(baseHandler)

export default handler