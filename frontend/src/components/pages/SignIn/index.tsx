import { VStack, Flex, Text, Heading} from '@chakra-ui/react'
import { Button } from '@components/shared/Button'
import { useAuth0 } from '@auth0/auth0-react'

const SignIn = () => {
    const { loginWithRedirect } = useAuth0()
    const onClickSignIn = () => loginWithRedirect()
    
    return (
        <Flex alignItems="center" justifyContent="center" mt="2rem">
            <VStack spacing="1rem">
                <Heading>Sign In</Heading>
                <Text>{`To access all the services you need to login first. If not, you can't move forward`}</Text>
                <Button w="200px" onClick={onClickSignIn}>Login</Button>
            </VStack>
        </Flex>
    )
}

export default SignIn