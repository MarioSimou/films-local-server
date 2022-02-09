import React from 'react'
import {Flex, VStack, Img, Heading, Grid, Link as ChakraLink, useToast, Spinner} from '@chakra-ui/react'
import type {NextPage} from 'next'
import NextLink from 'next/link'
import type { Song } from '@types'
import {format} from 'date-fns'
import {ChakraButton} from '@components/shared/Button'
import EditForm from '@components/pages/Song/components/EditForm'
import Label from '@components/pages/Song/components/Label'
import { useAuth0 } from '@auth0/auth0-react'
import { useRouter } from 'next/router'
import { useSong } from '@hooks'
import { fetchFile } from '@utils'

const SongPage: NextPage = () => {
    const [song, setSongs] = React.useState<Song | undefined>(undefined)
    const { isAuthenticated, isLoading: isAuthLoading } = useAuth0()

    const {deleteSong, isLoading, getSong} = useSong()
    const toast = useToast({
        isClosable: true,
        title: 'Delete',
        status: 'error',
        position: 'bottom-right',
    })
    const router = useRouter()
    const [isEdit, setIsEdit] = React.useState(false)
    const imageRef = React.useRef<File>()
    const hrefRef = React.useRef<File>()
    const enableEdit = React.useCallback(() => {
        return setIsEdit(() => true)
    }, [setIsEdit])
    const disableEdit = React.useCallback(() => {
        return setIsEdit(() => false)
    }, [setIsEdit])

    const onClickDeleteSong = React.useCallback(async () => {
        if(!song) {
            return 
        }

        const [e] = await deleteSong(song.guid)
        if(e){
            return toast({description: e.message})
        }
        toast({description: 'The song has been succesfully deleted', status: 'success'})
        return router.push('/')
    },[router, toast, song, deleteSong])

    React.useEffect(() => {
        const fetchSongAndImages = async () => {
            const guid = router.query.guid as string | undefined

            if(!guid) {
                return
            }
            const [e, song] = await getSong(guid)

            if(e){
                return toast({description: e.message})
            }
            if(!song){
                return toast({description: 'Song not found'})
            }

            const [imageError, imageFile] = await fetchFile("image",song.image)
            if(imageError){
                return toast({description: imageError.message})
            }

            const [hrefError, hrefFile] = await fetchFile("href",song.href)
            if(hrefError){
                return toast({description: hrefError.message})
            }

            imageRef.current = imageFile
            hrefRef.current = hrefFile

            setSongs(() => song)
        }
        fetchSongAndImages()
        return () => setSongs(undefined)
    }, [router.query.guid, setSongs, getSong])

    React.useEffect(() => {
        if(!isAuthenticated && !isAuthLoading){
            router.push('/sign-in')
        }
    }, [isAuthenticated, router, isAuthLoading])
    
    return (
        <VStack bg="gray.100" w="100%" minH="calc( 100vh - 84px)" p="2rem 0">
            {!song && <Spinner size="lg"/>}
            {song && <VStack alignItems="flex-start">
            <EditForm isOpen={isEdit} onClose={disableEdit} song={{...song, image: imageRef.current as File, href: hrefRef.current as File}}/>
            <Grid templateColumns={["1fr","1fr","1fr 1fr","1fr 1fr"]} gridGap="4rem" mt="2rem" bg="gray.50" p="4rem 2rem" w="100%">
                <Img src={song.image} boxSize="100%"/>
                <VStack spacing="2rem">
                    <VStack alignItems="flex-start" w="100%">
                        <Heading textTransform="capitalize">{song.name}</Heading>
                        {song.description && <Label title="Description" description={song.description}/>}
                        {song.href && <Label title="Href">
                            <NextLink href={song.href as string} passHref>
                                <ChakraLink color="blue.500">{song.href}</ChakraLink>
                            </NextLink>
                        </Label>}
                        <Label title="Created At" description={format(new Date(song.createdAt), 'EEEE, LLLL Lo, yyyy')}/>
                        <VStack w="100%" p="1rem 0 0.5rem 0">
                            <audio controls>
                                <source src={URL.createObjectURL(hrefRef.current as any)} type={hrefRef.current?.type}/>
                            </audio>
                        </VStack>
                    </VStack>
                    <Flex w="100%" gridGap="1rem">
                        <ChakraButton w="100%" colorScheme="yellow" onClick={enableEdit}>Edit</ChakraButton>
                        <ChakraButton w="100%" colorScheme="red" onClick={onClickDeleteSong} isLoading={isLoading}>Delete</ChakraButton>
                    </Flex>
                </VStack>
            </Grid>
        </VStack>}
        </VStack>
    )
}

export default SongPage