import axios from 'axios'

export const httpClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_BACKEND_API_HOST,
    timeout: 10000,
    headers: {
        Accept: 'application/json'
    }
})