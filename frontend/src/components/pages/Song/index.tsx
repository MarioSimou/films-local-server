import React from 'react'
import {Flex, VStack, Img, Heading, Grid, Link as ChakraLink, useToast} from '@chakra-ui/react'
import type {GetStaticProps, GetStaticPropsContext, GetStaticPaths, NextPage} from 'next'
import NextLink from 'next/link'
import type { Song as SongT, HTTPResponse, Song } from '@types'
import {httpClient} from '@utils'
import type { ParsedUrlQuery } from 'querystring'
import {format} from 'date-fns'
import {ChakraButton} from '@components/shared/Button'
import EditForm from '@components/pages/Song/components/EditForm'
import Label from '@components/pages/Song/components/Label'
import { useRouter } from 'next/router'
import type {AxiosResponse } from 'axios'
import { useSong } from '@hooks'

export type Props = {
    song: SongT
}

const Song: NextPage<Props> = ({song}) => {
    const {deletSong, isLoading} = useSong()
    const toast = useToast({
        isClosable: true,
        title: 'Delete',
        status: 'error',
        position: 'bottom-right',
    })
    const router = useRouter()
    const [isEdit, setIsEdit] = React.useState(false)
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

    return (
        <VStack alignItems="flex-start" >
            <EditForm isOpen={isEdit} onClose={disableEdit} song={song}/>
            <Grid templateColumns={["1fr","1fr","1fr 1fr","1fr 1fr"]} gridGap="4rem" mt="2rem" bg="gray.50" p="4rem 2rem">
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
    const {data}: AxiosResponse<HTTPResponse<SongT[]>> = await  httpClient.get('/api/v1/songs')

    if(!data.success){
        throw new Error(data.message)
    }

    const paths = data.data?.map((song: SongT) => ({params: {guid: song.guid}}))

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

    const {data}: AxiosResponse<HTTPResponse<SongT>> = await  httpClient.get(`/api/v1/songs/${guid}`)

    if(!data.success){
        throw new Error(data.message)
    }

    return {
        props: {
            song: data.data,
        },
        revalidate: true
    }
}

export default Song