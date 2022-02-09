import {decode, verify, Jwt} from 'jsonwebtoken'
import jwksRsa from 'jwks-rsa'
import { promisify } from 'util'
import type { APIGatewayProxyRequest, APIGatewayProxyEventPathParameters, APIGatewayProxyEventQueryStringParameters, Result } from '@libs/types'

const AUTH0_DOMAIN = process.env.AUTH0_DOMAIN


if(!AUTH0_DOMAIN){
    throw new Error("error: AUTH0 environment variables are not configured correctly")
}

const jwksClient = jwksRsa({
    jwksUri: `https://${AUTH0_DOMAIN}/.well-known/jwks.json`,
    cache: true,
    rateLimit: true,
    jwksRequestsPerMinute: 10,
})

const getSigningKey = promisify(jwksClient.getSigningKey)

const ErrTokenNotFound = new Error('error: token not found')
const ErrPayloadNotFound = new Error('error: payload not found')

export const getToken = <S = unknown, P = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters>(event: APIGatewayProxyRequest<S,P,Q>): Result<string> => {
    const authorization = event.headers['Authorization'] ?? event.headers['authorization']
    if(!authorization){
        return [ErrTokenNotFound]
    }

    const [,token] = authorization.match(/^Bearer (.*)$/)
    if(!token){
        return [ErrTokenNotFound]
    }
    return [undefined, token]
}

export const getJWKsPublicKey = async (token: string): Promise<Result<string>> => {
    const payload = decode(token, {complete: true})

    if(!payload){
        return [ErrPayloadNotFound]
    }

    const options = await getSigningKey(payload.header.kid)
    const publicKey = options.getPublicKey()
    return [undefined, publicKey]
}

export const verifyToken = async <S = unknown, P = APIGatewayProxyEventPathParameters, Q = APIGatewayProxyEventQueryStringParameters>(event: APIGatewayProxyRequest<S,P,Q>): Promise<Result<Jwt>> => {
    try {
        const [getTokenError, jwksToken] = getToken(event)
        if(getTokenError){
            return [getTokenError]
        }
        
        const [getJWKSPublicKeyError, publicKey] = await getJWKsPublicKey(jwksToken)
        if(getJWKSPublicKeyError){
            return [getJWKSPublicKeyError]
        }

        const jwtToken= await verify(jwksToken, publicKey, {complete: true}) as Jwt
        return [undefined, jwtToken]
    }catch(e){
        console.log("error: ", e)
        return [e]
    }
}