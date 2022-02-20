import {APIGatewayProxyEvent, APIGatewayProxyResult} from 'aws-lambda'

export type Request<
  B extends Object = unknown,
  P extends APIGatewayProxyEvent['pathParameters'] = {},
  Q extends APIGatewayProxyEvent['queryStringParameters'] = {}
> = Omit<APIGatewayProxyEvent, 'body' | 'pathParameters' | ''> & {
  body: B
  pathParameters: P
  queryStringParameters: Q
}
export type Response = APIGatewayProxyResult

export type Handler<
  B extends Object = unknown,
  P extends APIGatewayProxyEvent['pathParameters'] = {},
  Q extends APIGatewayProxyEvent['queryStringParameters'] = {}
> = (event: Request<B, P, Q>) => Promise<Response>

export type Status = 'pending' | 'completed'
export type GUID = string
export type Song = {
  guid: GUID
  title: string
  description?: string
  image?: string
  imageStatus?: Status
  href?: string
  hrefStatus?: Status
  createdAt: string
}

export type Result<R = unknown> = [Error] | [undefined, R]
