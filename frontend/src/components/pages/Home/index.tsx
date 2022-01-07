import {Grid, VStack, Heading, Collapse, Text} from '@chakra-ui/react'
import type {GetStaticProps, NextPage} from 'next'
import type {AxiosResponse} from 'axios'
import type { Song, HTTPResponse } from '@types'
import {httpClient, getBackendURL, isHTTPError} from '@utils'
import SongCard from '@components/pages/Home/components/SongCard'

export type Props = {
    songs: Song[]
}

const Home: NextPage<Props> = ({songs}) => {
    const isSongs = songs.length > 0

    return (
        <VStack bg="gray.100" w="100%" minH="calc( 100vh - 84px)">
            <VStack p="2rem" w="100%" maxW="1200px" alignItems="flex-start">
                <Heading mb="1rem" fontSize="1.75rem">Your songs</Heading>
                <Collapse in={!isSongs}>
                    <Text>Songs not found. Please insert one</Text>
                </Collapse>
                <Collapse in={isSongs}>
                    <Grid templateColumns="repeat(auto-fill, minmax(250px, 1fr))" gridGap="1rem">
                        {songs.map(song => {
                            return (
                                <SongCard key={song.guid} song={song}/>
                            )
                        })}
                    </Grid>
                </Collapse>
            </VStack>
        </VStack>
    )
}

export const getStaticProps: GetStaticProps = async () => {
    try {
        const {data}: AxiosResponse<HTTPResponse<Song[]>> = await  httpClient({
            url: getBackendURL('/api/v1/songs'),
            method: 'GET',
        })

        if(!data.success){
            throw new Error(data.message)
        }
    
        const songs = data.status === 404 ? [] : data.data
    
        return {
            props: {
                songs,
            },
            revalidate: 1
        }
    }catch(e: any){
        if(isHTTPError(e?.response?.data)){
            if(e.response.data.status === 404){
                return {
                    props: {
                        songs: []
                    },
                    revalidate: 1,        
                }
            }
        }

        throw e
    }

}

export default Home