import React from 'react'
import { Box, VStack, List, ListItem, Text, CloseButton } from '@chakra-ui/react'
import {useRouter} from 'next/router'

export type Props = {isOpen: boolean, onClose: () => void}

const SidebarItem = ({onClick, content}: {onClick: () => void, content: string}) => {
    return (
        <ListItem w="100%" fontSize="1.25rem" _hover={{bg: "primary.700", cursor: "pointer"}} p="0.75rem 2rem" onClick={onClick}>
        <Text>{content}</Text>
    </ListItem>
    )
}

const Sidebar = ({isOpen, onClose}: Props) => {
    const router = useRouter()
    const left = isOpen ? "0vw": "-100vw"
    const onClickListItem = React.useCallback((href: string) => {
        return () => {
            onClose()
            router.push(href)
        }
    }, [onClose, router])

    return (
        <React.Fragment>
            <Box visibility={isOpen ? "visible" : "hidden"} position="fixed" top="0" left="0" w="100vw" h="100vh" bg={isOpen ? "blackAlpha.600": "transparent"} transition="background ease-in 0.3s" onClick={onClose}/>
            <VStack zIndex="10" transition="left ease-in-out 0.4s" h="100vh" bg="primary.800" position="fixed" left={left} width={["100%","100%", "300px", "300px"]} alignItems="flex-start">
                <CloseButton color="text.200" onClick={onClose} alignSelf="flex-end" size="lg" _hover={{bg: "transparent"}} _active={{bg: "transparent"}} m="1.75rem 1.5rem"/>                    
                <List color="text.200" letterSpacing="wide" w="100%">
                    <SidebarItem content="Home" onClick={onClickListItem('/')}/>
                    <SidebarItem content="Upload Song" onClick={onClickListItem('/upload')}/>
                </List>    
            </VStack>
        </React.Fragment>
    )
}

export default Sidebar