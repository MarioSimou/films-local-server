import {Flex, Img, VStack, Heading, Text} from '@chakra-ui/react'
import type { Song } from '@types'
import NextLink from 'next/link'

export type Props = {
    song: Song
}

const SongCard = ({song}: Props) => {
    return (
        <NextLink href={`/songs/${song.guid}`}>
            <Flex as="a" flexDirection="column" _hover={{cursor: 'pointer'}}>
                <Img src={song.image} alt={song.name} boxSize="100%"/>
                <Flex flexDirection="column" bg="gray.50" p="0.5rem 1rem" >
                    <Heading fontSize="1.5rem" as="h3" textTransform="capitalize">{song.name}</Heading>
                    <Text color="gray.600">{song.description}</Text>
                </Flex>
            </Flex>
        </NextLink>
    )
}

export default SongCard