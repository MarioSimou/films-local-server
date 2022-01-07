import React from 'react'
import {Flex, VStack, Img, Heading, Grid, Link as ChakraLink, useToast, Collapse, Spinner} from '@chakra-ui/react'
import type {GetStaticProps, GetStaticPropsContext, GetStaticPaths, NextPage} from 'next'
import NextLink from 'next/link'
import type { Song, HTTPResponse } from '@types'
import {httpClient} from '@utils'
import type { ParsedUrlQuery } from 'querystring'
import {format} from 'date-fns'
import {ChakraButton} from '@components/shared/Button'
import EditForm from '@components/pages/Song/components/EditForm'
import Label from '@components/pages/Song/components/Label'
import { useRouter } from 'next/router'
import type {AxiosResponse } from 'axios'
import { useSong } from '@hooks'
import { fetchFile, getBackendURL } from '@utils'

export type Props = {
    song: Song
}

const fallbackSong: Song = {
    guid: '',
    name: '',
    description: '',
    image: '',
    location: '',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
}

const SongPage: NextPage<Props> = ({song = fallbackSong}) => {
    const [areFilesLoaded, setAreFilesLoaded] = React.useState(false)
    const {deletSong, isLoading} = useSong()
    const toast = useToast({
        isClosable: true,
        title: 'Delete',
        status: 'error',
        position: 'bottom-right',
    })
    const router = useRouter()
    const [isEdit, setIsEdit] = React.useState(false)
    const imageRef = React.useRef<File>()
    const locationRef = React.useRef<File>()
    const enableEdit = React.useCallback(() => {
        return setIsEdit(() => true)
    }, [setIsEdit])
    const disableEdit = React.useCallback(() => {
        return setIsEdit(() => false)
    }, [setIsEdit])

    const onClickDeleteSong = React.useCallback(async () => {
        const [e] = await deletSong(song.guid)
        if(e){
            return toast({description: e.message})
        }
        toast({description: 'The song has been succesfully deleted', status: 'success'})
        return router.push('/')
    },[router, toast, song.guid, deletSong])

    React.useEffect(() => {
        const fetchFiles = async () => {
            const [imageError, imageFile] = await fetchFile("image",song.image)
            if(imageError){
                return toast({description: imageError.message})
            }

            const [locationError, locationFile] = await fetchFile("location",song.location)
            if(locationError){
                return toast({description: locationError.message})
            }

            imageRef.current = imageFile
            locationRef.current = locationFile
            setAreFilesLoaded(() => true)
        }
        fetchFiles()
    }, [])

    if(router.isFallback){
        return (<Spinner size="lg"/>)
    }

    return (
        <VStack alignItems="flex-start" bg="gray.100" w="100%" minH="calc( 100vh - 84px)">
            {areFilesLoaded && <EditForm isOpen={isEdit} onClose={disableEdit} song={{...song, image: imageRef.current as File, location: locationRef.current as File}}/>}
            <Grid templateColumns={["1fr","1fr","1fr 1fr","1fr 1fr"]} gridGap="4rem" mt="2rem" bg="gray.50" p="4rem 2rem" w="100%">
                <Img src={song.image} boxSize="100%"/>
                <VStack spacing="2rem">
                    <VStack alignItems="flex-start" w="100%">
                        <Heading textTransform="capitalize">{song.name}</Heading>
                        {song.description && <Label title="Description" description={song.description}/>}
                        <Label title="Location">
                            <NextLink href={song.location as string} passHref>
                                <ChakraLink color="blue.500">{song.location}</ChakraLink>
                            </NextLink>
                        </Label>
                        <Label title="Created At" description={format(new Date(song.createdAt), 'EEEE, LLLL Lo, yyyy')}/>
                        <Label title="Updated At" description={format(new Date(song.updatedAt), 'EEEE, LLLL Lo, yyyy')}/>
                        <VStack w="100%" p="1rem 0 0.5rem 0">
                            {!areFilesLoaded && <Spinner/>}
                            {areFilesLoaded && <audio controls>
                                <source src={URL.createObjectURL(locationRef.current as any)} type={locationRef.current?.type}/>
                            </audio>}
                        </VStack>
                    </VStack>
                    <Flex w="100%" gridGap="1rem">
                        <ChakraButton w="100%" colorScheme="yellow" onClick={enableEdit}>Edit</ChakraButton>
                        <ChakraButton w="100%" colorScheme="red" onClick={onClickDeleteSong} isLoading={isLoading}>Delete</ChakraButton>
                    </Flex>
                </VStack>
            </Grid>
        </VStack>
    )
}

export const getStaticPaths: GetStaticPaths = async () => {
    try {
        const {data} : AxiosResponse<HTTPResponse<Song[]>> = await  httpClient({
            url: getBackendURL('/api/v1/songs'),
            method: 'GET',
        })
        const paths = data.data?.map((song: Song) => ({params: {guid: song.guid}}))
    
        if(!paths) {
            return {
                paths: [],
                fallback: true,
            }
        }
    
        return {
            paths,
            fallback: true,
        }
    } catch(e: any) {
        return {
            paths: [],
            fallback: true,
        }
    }
}

export const getStaticProps: GetStaticProps = async ({params}: GetStaticPropsContext) => {
    const {guid} = params as ParsedUrlQuery

    if(!guid){
        return {
            redirect: {
                permanent: false,
                destination: '/'
            }
        }
    }

    const {data}: AxiosResponse<HTTPResponse<Song>> = await  httpClient.get(getBackendURL(`/api/v1/songs/${guid}`))

    if(!data.success){
        throw new Error(data.message)
    }

    return {
        props: {
            song: data.data,
        },
        revalidate: 1
    }
}

export default SongPage