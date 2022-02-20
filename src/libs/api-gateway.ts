import middy from '@middy/core'
import middyJsonBodyParser from '@middy/http-json-body-parser'
import type {MiddyfiedHandler} from '@middy/core'
import type {Handler, Response, Request, Result} from '@libs/types'
import {ValidationError} from 'joi'
import {parse} from 'aws-multipart-parser'
import {FileData, FormData} from 'aws-multipart-parser/dist/models'

export const ErrUnsupportedMediaType = new Error('error: unsupported media type')

export const formatJSONResponse = (response: Record<string, unknown>) => {
  return {
    statusCode: 200,
    body: JSON.stringify(response),
  }
}

const isFileData = (options: any): options is FileData => options instanceof Object

export const parseMultipartFormData = (event: Request<string>): Result<FileData> => {
  const contentType = event.headers['Content-Type'] ?? event.headers['content-type']
  const isMultipartFormData = /multipart\/form-data/.test(contentType)
  if (!isMultipartFormData) {
    return [ErrUnsupportedMediaType]
  }

  const result = parse(event, true)
  const [options] = Object.values<FormData | undefined>(result)

  if (!options) {
    const e = new Error('error: missing form data')
    return [e]
  }

  if (!isFileData(options)) {
    const e = new Error('error: invalid field data')
    return [e]
  }

  return [undefined, options]
}

export const newHandler = (handler: Handler): MiddyfiedHandler => {
  return middy(handler).use(middyJsonBodyParser())
}

export type StatusCode = number
export const StatusOk: StatusCode = 200
export const StatusNoContent: StatusCode = 204
export const StatusCreated: StatusCode = 201
export const StatusBadRequest: StatusCode = 400
export const StatusUnauthorized: StatusCode = 401
export const StatusForbidden: StatusCode = 403
export const StatusNotFound: StatusCode = 404
export const StatusRequestTimeout: StatusCode = 408
export const StatusUnsupportedMediaType = 415
export const StatusInternalServerError: StatusCode = 500

type ResponseBody<B = unknown> = {
  status: StatusCode
  success: boolean
  message?: string
  body?: B
}

const isError = (e: any): e is Error => {
  return e instanceof Error || e instanceof ValidationError
}

const newResponseBody = <B = unknown>(statusCode: Response['statusCode'], body: any): ResponseBody<B> => {
  switch (statusCode) {
    case StatusOk:
    case StatusCreated: {
      return {
        status: statusCode,
        success: true,
        body: body,
      }
    }
    case StatusNoContent: {
      return {
        status: statusCode,
        success: true,
      }
    }
    default:
      if (isError(body)) {
        return {
          status: statusCode,
          success: false,
          message: body.message,
        }
      }

      return {
        status: statusCode,
        success: false,
        message: body as string,
      }
  }
}

export const newResponse = (
  statusCode: Response['statusCode'],
  body?: unknown,
  headers: Record<string, string> = {}
): Response => {
  return {
    statusCode,
    body: JSON.stringify(newResponseBody(statusCode, body)),
    headers,
  }
}
