import React from 'react'
import { Flex, Box, Icon } from '@chakra-ui/react'
import { GiHamburgerMenu as MenuIcon} from 'react-icons/gi'
import Sidebar from '@components/shared/Navbar/components/Sidebar'

export type Props = {}

const Navbar: React.FC<Props> = ({children}) => {
    const [isOpen, setIsOpen] = React.useState<boolean>(false)
    const toggleSidebar = React.useCallback(() => setIsOpen(open => !open), [setIsOpen]) 

    return (
        <Flex flexDirection="column">
            <Flex as="header" p="2rem" justifyContent="flex-end" bg="primary.900">
                <Icon as={MenuIcon} _hover={{cursor: 'pointer'}} w={5} h={5} color="secondary.500" onClick={toggleSidebar}/>
            </Flex>
            <Sidebar isOpen={isOpen} onClose={toggleSidebar}/>
            <Box as="main">
                {children}
            </Box>
        </Flex>
    )
}

export default Navbar