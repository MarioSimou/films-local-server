import middy from '@middy/core'
import { verifyToken } from '@libs/jwt'
import type { APIGatewayProxyRequest, APIGatewayProxyEventPathParameters, APIGatewayProxyEventQueryStringParameters, APIGatewayProxyResponse } from '@libs/types'
import { formatJSONResponse, StatusUnauthorized, ErrUnauthorized } from '@libs/api-gateway'

const authMiddleware = <S = unknown, P = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters>(): middy.MiddlewareObj<APIGatewayProxyRequest<S,P,Q>, APIGatewayProxyResponse> => {
    const before: middy.MiddlewareFn<APIGatewayProxyRequest<S,P,Q>, APIGatewayProxyResponse> = async ({event}): Promise<APIGatewayProxyResponse | undefined> => {
        const [e, jwt] = await verifyToken(event) 
        if(e){
            return formatJSONResponse(StatusUnauthorized, ErrUnauthorized.message)
        }
        event.jwt = jwt
    } 

    return {
        before,
    }
}

export default authMiddleware