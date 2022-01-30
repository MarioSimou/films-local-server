import React from 'react'
import {httpClient, isHTTPError, getBackendURL} from '@utils'
import type {AxiosResponse} from 'axios'
import type { Song, PostSong, PutSong, HTTPResponse } from '@types'

type APIGatewayResponse<T = unknown> = [Error] | [undefined, T]

type GetSongsResponse = APIGatewayResponse<Song[]>
type GetSongResponse = APIGatewayResponse<Song>
type DeleteSongResponse = [Error | undefined] 
type PostSongResponse = APIGatewayResponse<Song>

export const useSong = () => {
    const [isLoading, setIsLoading] = React.useState(false)

    const getSongs = React.useCallback(async (): Promise<GetSongsResponse> => {
        try {
            setIsLoading(() => true)

            const {data}: AxiosResponse<HTTPResponse<Song[]>> = await httpClient({
                url: getBackendURL('/api/v1/songs'),
                method: 'GET'
            })
            if(!data.success){
                const e = new Error(data.message)
                return [e]
            }
            setIsLoading(() => false)
            return [undefined, data.data as Song[]]
        } catch(e: any){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [new Error(e.response.data.message)]
            }
            return [e]
        }
    }, [])

    const getSong = React.useCallback(async (guid: string): Promise<GetSongResponse> => {
        try {
            setIsLoading(() => true)

            const {data}: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                url: getBackendURL(`/api/v1/songs/${guid}`),
                method: 'GET'
            })
            if(!data.success){
                const e = new Error(data.message)
                return [e]
            }
            setIsLoading(() => false)
            return [undefined, data.data as Song]
        } catch(e: any){
            if(isHTTPError(e?.response?.data)){
                return [new Error(e.response.data.message)]
            }
            setIsLoading(() => false)
            return [e]
        }
    }, [])

    const deletSong = React.useCallback(async (songGUID: string): Promise<DeleteSongResponse>  => {
        try {
            setIsLoading(() => true)
            await httpClient({
                url: getBackendURL(`/api/v1/songs/${songGUID}`),
                method: 'DELETE',
            })
            setIsLoading(() => false)
            return [undefined]
        }catch(e: any){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [new Error(e.response.data.message)]
            }
            return [e]
        }
    }, [])

    const postSong = React.useCallback(async (song: PostSong): Promise<PostSongResponse> => {
        try{
            setIsLoading(true)
            const {name, description, image, href} = song 
            const postSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                url: getBackendURL('/api/v1/songs'),
                method: 'POST',
                data: JSON.stringify({name, description})
            })
    
            if(!postSongResponse.data.success){
                throw new Error(postSongResponse.data.message)
            }
    
            const songGUID = postSongResponse.data.data?.guid
            const putImageFormData = new FormData()
            putImageFormData.append('image', image)
            const putSongImageResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                method: 'PUT',
                url: getBackendURL(`/api/v1/songs/${songGUID}/image`), 
                data: putImageFormData
            })
    
            if(!putSongImageResponse.data.success){
                throw new Error(putSongImageResponse.data.message)
            }
    
            const uploadSongFormData = new FormData()
            uploadSongFormData.append('song', href)
            const uploadSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                method: 'PUT',
                url: getBackendURL(`/api/v1/songs/${songGUID}/upload`),
                data: uploadSongFormData,
            })
    
            if(!uploadSongResponse.data.success){
                throw new Error(uploadSongResponse.data.message)
            }
            setIsLoading(false)
            return [undefined, uploadSongResponse.data.data as Song]
        }catch(e: any){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [new Error(e.response.data.message)]
            }
            return [e]
        }
    }, [setIsLoading])

    const putSong = React.useCallback(async (songGUID: string, song: PutSong): Promise<PostSongResponse> => {
        try{
            setIsLoading(true)
            const {name, description, image, href} = song 
            const putSongPayload = JSON.stringify({name, description})
            const postSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                method: 'PUT',
                url: getBackendURL(`/api/v1/songs/${songGUID}`), 
                data: putSongPayload
            })
    
            if(!postSongResponse.data.success){
                throw new Error(postSongResponse.data.message)
            }
    
            const putImageFormData = new FormData()
            putImageFormData.append('image', image)
            const putSongImageResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                method: 'PUT',
                url: getBackendURL(`/api/v1/songs/${songGUID}/image`), 
                data: putImageFormData
            })
    
            if(!putSongImageResponse.data.success){
                throw new Error(putSongImageResponse.data.message)
            }
    
            const uploadSongFormData = new FormData()
            uploadSongFormData.append('song', href)
            const uploadSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient({
                method: 'PUT',
                url: getBackendURL(`/api/v1/songs/${songGUID}/upload`), 
                data: uploadSongFormData
            })
    
            if(!uploadSongResponse.data.success){
                throw new Error(uploadSongResponse.data.message)
            }
            setIsLoading(false)
            return [undefined, uploadSongResponse.data.data as Song]
        }catch(e: any){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [new Error(e.response.data.message)]
            }
            return [e]
        }
    }, [setIsLoading])

    return {
        getSong,
        getSongs,
        putSong,
        postSong,
        isLoading,
        deletSong
    }
}