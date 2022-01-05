import {Grid, VStack, Heading, Collapse, Text} from '@chakra-ui/react'
import type {GetStaticProps, NextPage} from 'next'
import type {AxiosResponse} from 'axios'
import type { Song, HTTPResponse } from '@types'
import {httpClient, isHTTPError} from '@utils'
import SongCard from '@components/pages/Home/components/SongCard'

export type Props = {
    songs: Song[]
}

const Home: NextPage<Props> = ({songs}) => {
    const isSongs = songs.length > 0

    return (
        <VStack p="2rem" maxW="1200px" alignItems="flex-start" >
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
    )
}

export const getStaticProps: GetStaticProps = async () => {
    try {
        const {data}: AxiosResponse<HTTPResponse<Song[]>> = await  httpClient.get('/api/v1/songs')

        if(!data.success){
            throw new Error(data.message)
        }
    
        const songs = data.status === 404 ? [] : data.data
    
        return {
            props: {
                songs,
            },
            revalidate: true
        }
    }catch(e){
        if(isHTTPError(e?.response?.data)){
            if(e.response.data.status === 404){
                return {
                    props: {
                        songs: []
                    },
                    revalidate: true,        
                }
            }
        }

        throw e
    }

}

export default Home