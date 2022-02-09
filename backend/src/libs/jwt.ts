import {decode, verify, Jwt} from 'jsonwebtoken'
import jwksRsa from 'jwks-rsa'
import { promisify } from 'util'
import type { Result, AuthEvent } from '@libs/types'

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

export const getToken = ({authorizationToken}: AuthEvent): Result<string> => {
    if(!authorizationToken){
        return [ErrTokenNotFound]
    }

    const [,token] = authorizationToken.match(/^Bearer (.*)$/)
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

export const verifyToken = async (event: AuthEvent): Promise<Result<Jwt>> => {
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
        return [e]
    }
}