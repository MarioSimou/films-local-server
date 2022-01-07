import React from 'react'
import {httpClient, isHTTPError, getBackendURL} from '@utils'
import type {AxiosResponse} from 'axios'
import type { Song, PostSong, PutSong, HTTPResponse } from '@types'

type DeleteSongResponse = [Error | undefined]
type PostSongResponse = [Error] | [undefined, Song]

export const useSong = () => {
    const [isLoading, setIsLoading] = React.useState(false)
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
            const {name, description, image, location} = song 
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
            uploadSongFormData.append('song', location)
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
            console.log(e.message)
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
            const {name, description, image, location} = song 
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
            uploadSongFormData.append('song', location)
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
        putSong,
        postSong,
        isLoading,
        deletSong
    }
}