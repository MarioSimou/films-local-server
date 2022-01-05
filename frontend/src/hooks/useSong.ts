import React from 'react'
import {httpClient, isHTTPError} from '@utils'
import type {AxiosResponse} from 'axios'
import type { Song, PostSong, HTTPResponse } from '@types'

type DeleteSongResponse = [Error | undefined]
type PostSongResponse = [Error] | [undefined, Song]

export const useSong = () => {
    const [isLoading, setIsLoading] = React.useState(false)
    const deletSong = React.useCallback(async (songGUID: string): Promise<DeleteSongResponse>  => {
        try {
            setIsLoading(() => true)
            await httpClient.delete(`/api/v1/songs/${songGUID}`)
            setIsLoading(() => false)
            return [undefined]
        }catch(e){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [e.response.data.message]
            }
            return [e.message]
        }
    }, [])

    const postSong = React.useCallback(async (song: PostSong): Promise<PostSongResponse> => {
        try{
            setIsLoading(true)
            const {name, description, image, location} = song 
            const postSongPayload = JSON.stringify({name, description})
            const postSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient.post('/api/v1/songs', postSongPayload)
    
    
            if(!postSongResponse.data.success){
                throw new Error(postSongResponse.data.message)
            }
    
            const songGUID = postSongResponse.data.data?.guid
            const putImageFormData = new FormData()
            putImageFormData.append('image', image)
            const putSongImageResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient.put(`/api/v1/songs/${songGUID}/image`, putImageFormData)
    
            if(!putSongImageResponse.data.success){
                throw new Error(putSongImageResponse.data.message)
            }
    
            const uploadSongFormData = new FormData()
            uploadSongFormData.append('song', location)
            const uploadSongResponse: AxiosResponse<HTTPResponse<Song>> = await httpClient.put(`/api/v1/songs/${songGUID}/upload`, putImageFormData)
    
            if(!uploadSongResponse.data.success){
                throw new Error(uploadSongResponse.data.message)
            }
            setIsLoading(false)
            return [undefined, uploadSongResponse.data.data as Song]
        }catch(e){
            setIsLoading(() => false)
            if(isHTTPError(e?.response?.data)){
                return [e.response.data.message]
            }
            return [e.message]
        }
    }, [setIsLoading])

    return {
        postSong,
        isLoading,
        deletSong
    }
}