import React from 'react'
import {Flex, Img, Heading, Text} from '@chakra-ui/react'
import type { Song } from '@types'
import { useRouter } from 'next/router'

export type Props = {
    song: Song
}

const SongCard = ({song}: Props) => {
    const router = useRouter()

    const onClickCard = React.useCallback(() => {
        return router.push(`/songs/${song.guid}`)
    }, [song.guid, router])

    return (
        <Flex as="a" flexDirection="column" _hover={{cursor: 'pointer'}} onClick={onClickCard}>
        <Img src={song.image} alt={song.name} boxSize="100%"/>
        <Flex flexDirection="column" bg="gray.50" p="0.5rem 1rem" >
            <Heading fontSize="1.5rem" as="h3" textTransform="capitalize">{song.name}</Heading>
            <Text color="gray.600">{song.description}</Text>
        </Flex>
    </Flex>
    )
}

export default SongCard