import type { Handler, APIGatewayProxyEventPathParameters, APIGatewayProxyEventQueryStringParameters, APIGatewayProxyResponse } from '@libs/types'
import type { APIGatewayProxyEvent } from 'aws-lambda'
import jsonBodyParser from '@middy/http-json-body-parser'
import middy from '@middy/core'
import type { MiddlewareObj } from '@middy/core'
import cors from '@middy/http-cors'
import { parse } from 'aws-multipart-parser'
import auth  from '@libs/middlewares/auth'

export const newHandler = <T, S = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters>(handler: Handler<T, S, Q>, ...middlewares: MiddlewareObj[]  ) => middy(handler).use([
    jsonBodyParser(),
    auth(),
    cors({
        origins: [ process.env.ALLOW_ORIGIN ],
        methods: 'GET,POST,PUT,DELETE,OPTIONS,PATCH',
        maxAge: 86400,
    }),
    ...middlewares,
])

export type StatusCode = number
export const StatusBadRequest: StatusCode = 400
export const StatusNotFound: StatusCode = 404
export const StatusForbidden: StatusCode = 403
export const StatusUnauthorized: StatusCode = 401
export const StatusInternalServerError: StatusCode = 500
export const StatusOk: StatusCode = 200
export const StatusNoContent: StatusCode = 204

export const ErrNotFound = new Error('error: not found')
export const ErrNoMultipart = new Error('error: content-type should be mutlipart/form-data')
export const ErrMultipartNotFound = new Error("error: multipart not found")
export const ErrUnauthorized = new Error('error: unauthorized')

type HTTPResponse<T> = {
    status: StatusCode
    success :boolean
    message?: string
    data?: T
}

const newResponse = <T>(status: StatusCode, data: T | string): HTTPResponse<T> => {
    switch(status){
        case StatusBadRequest:
        case StatusUnauthorized:
        case StatusForbidden:
        case StatusNotFound:
        case StatusInternalServerError: {
            return {
                status,
                success: false,
                message: data as string,
            }
        }
        case StatusOk:
        case StatusNoContent: {
            return {
                status,
                success: true,
                data: data as T,
            }
        }
        default: {
            return {
                status: StatusBadRequest,
                success: false,
                message: "invalid status code"
            }
        }
    }
}

export const formatJSONResponse = <T = string>(statusCode: StatusCode, body?: T): APIGatewayProxyResponse => {
    if(!body){
        return {statusCode, body: ""}
    }
  return {
    statusCode,
    body: JSON.stringify(newResponse<T>(statusCode, body)),
  }
}

export const isMultipart = (headers: Record<string, string>): boolean => {
    const contentType = headers['content-type'] ?? headers['Content-Type']
    const isMultipart = /multipart\/form-data/.test(contentType)
    return isMultipart
}

export type Multipart = {
    type: string
    filename: string
    contentType: string
    content: Buffer
}

type ParseMultipartResult = [Error] | [undefined, Multipart]
export const parseMultipart = (event: APIGatewayProxyEvent): ParseMultipartResult => {
    const res = parse(event, true)
    const [multipart] = Object.values(res)
    if(!multipart){
        return [ErrMultipartNotFound]
    }
    return [undefined, multipart as Multipart]
}