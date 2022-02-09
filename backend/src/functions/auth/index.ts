import { verifyToken } from '@libs/jwt'
import {AuthResponse, AuthEvent} from '@libs/types'

type Effect = 'Allow' | 'Deny'
const getPolicyDocument = (effect: Effect) => ({
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": effect,
            "Action": [
                "execute-api:Invoke",
            ],
            "Resource": "*",
        }
    ],
})

const auth = async (event: AuthEvent): Promise<AuthResponse> => {
    const [verifyTokenError] = await verifyToken(event)
    if(verifyTokenError){
        return {
            principalId: event.methodArn,
            policyDocument: getPolicyDocument('Deny')
        }
    }

    return {
        principalId: event.methodArn,
        policyDocument: getPolicyDocument('Allow'),
    }
}

export default auth