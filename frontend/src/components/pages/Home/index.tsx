import React from 'react'
import {Grid, VStack, Heading, Collapse, Text, useToast, Spinner} from '@chakra-ui/react'
import type {NextPage} from 'next'
import type { Song } from '@types'
import { useSong } from '@hooks'
import SongCard from '@components/pages/Home/components/SongCard'
import { useRouter } from 'next/router'
import { useAuth0 } from '@auth0/auth0-react'

const Home: NextPage = () => {
    const [songs, setSongs] = React.useState<Song[]>([])
    const {getSongs, isLoading} = useSong()
    const { isAuthenticated, isLoading: isAuthLoading } = useAuth0()
    const router = useRouter()
    const toast = useToast({
        position: 'bottom-right',
        isClosable: true,
        title: 'Home',
        status: 'error',
    })

    React.useEffect(() => {
        const fetchSongs = async () => {
            const [e, songs] = await getSongs()

            if(e && !/not found/.test(e.message)){
                return toast({description: e.message})
            }

            setSongs(() => songs ?? [])
        }
        fetchSongs()
        return () => setSongs([])
    }, [getSongs])

    React.useEffect(() => {
        if(!isAuthenticated && !isAuthLoading){
            router.push('/sign-in')
        }
    }, [isAuthenticated, isAuthLoading, router]) 

    return (
        <VStack bg="gray.100" w="100%" minH="calc( 100vh - 84px)">
            <VStack p="2rem" w="100%" maxW="1200px" alignItems="flex-start">
                <Heading mb="1rem" fontSize="1.75rem">Your songs</Heading>
                {isLoading && <Spinner size="lg" alignSelf="center"/>}
                <Collapse in={songs.length == 0 && !isLoading}>
                    <Text>Songs not found. Please insert one</Text>
                </Collapse>
                <Collapse in={songs.length !== 0}>
                    <Grid templateColumns="repeat(auto-fill, minmax(250px, 1fr))" gridGap="1rem">
                        {songs.map(song => {
                            return (<SongCard key={song.guid} song={song}/>)
                        })}
                    </Grid>
                </Collapse>
            </VStack>
        </VStack>
    )
}

export default Home