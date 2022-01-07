import axios from 'axios'

export const getBackendURL = (path: string): string => {
    if(process.env.NEXT_PUBLIC_BACKEND_STAGE){
        const stage = process.env.NEXT_PUBLIC_BACKEND_STAGE
        return `/${stage}${path}`
    }
    return path
}

export const httpClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_BACKEND_API_HOST,
    timeout: 10000,
    headers: {
        Accept: 'application/json'
    }
})