import { HTTPResponse } from '@types'
export * from './httpClient'

export const isHTTPError = (data: any): data is HTTPResponse<unknown> => {
    return data?.status && !data?.success && data?.message
} 
