import { HTTPResponse } from '@types'
export * from './httpClient'

export const isHTTPError = (data: any): data is HTTPResponse<unknown> => {
    return data?.status && !data?.success && data?.message
} 

export const fetchFile = async (name: string, href: string): Promise<[Error] | [undefined, File]> => {
    try {
        const res = await fetch(href)
        const blob = await res.blob()
        const file = new File([blob], name)
        return [undefined, file]
    } catch(e: any){
        return [e.message]
    }
}