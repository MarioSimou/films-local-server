import type {Jwt} from 'jsonwebtoken'
import type { APIGatewayProxyEvent, APIGatewayProxyResult, APIGatewayProxyEventPathParameters, APIGatewayProxyEventQueryStringParameters } from "aws-lambda"
export type {APIGatewayProxyEventPathParameters, APIGatewayProxyEventQueryStringParameters}
export type APIGatewayProxyRequest<S = unknown, P = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters> = Omit<APIGatewayProxyEvent, 'body' | 'pathParameters' | 'queryStringParameters'> & { body: S, pathParameters: P, queryStringParameters: Q | null, jwt: Jwt}
export type APIGatewayProxyResponse = APIGatewayProxyResult
export type Handler<T = unknown, P = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters> = (event: APIGatewayProxyRequest<T, P, Q>) => Promise<APIGatewayProxyResponse>

export type Timestamp = {
    createdAt: string
}
export type IGUID = {
    guid: string
}

export type Song = {
    name: string
    description?: string
    image?: string
    href?: string
} & Timestamp & IGUID
 
export type HTTPResponse<T> = [Error] | [undefined, T]
export type Result<T> = HTTPResponse<T>
export type GUIDPathParameters  = { guid: string }