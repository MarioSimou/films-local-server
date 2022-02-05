import { formatJSONResponse, newHandler, StatusOk, StatusInternalServerError, StatusNotFound, ErrNotFound } from '@libs/api-gateway'
import { Handler, APIGatewayProxyEventPathParameters } from '@libs/types'
import { findAllSongs } from '@libs/dynamodb'

const baseHandler: Handler<unknown, APIGatewayProxyEventPathParameters, {limit?: string}> = async (event) => {
    const { limit } = event.queryStringParameters ?? { limit: undefined }

    const [e,songs] = await findAllSongs(parseInt(limit))
    if(e){
        return formatJSONResponse(StatusInternalServerError, e.message)
    }
    if(songs.length === 0) {
        return formatJSONResponse(StatusNotFound, ErrNotFound.message)
    }
    return formatJSONResponse(StatusOk, songs)
}

const handler = newHandler(baseHandler)

export default handler