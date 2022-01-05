import {Button as ChakraButton, ButtonProps} from '@chakra-ui/react'
export {Button as ChakraButton} from '@chakra-ui/react'

export const Button = ({children, ...other}: ButtonProps) => {
    return (                        
        <ChakraButton isFullWidth bg="primary.900" color="white" _hover={{opacity: '0.8'}} {...other}>{children}</ChakraButton>
    )
}

export default Button