import React from "react"
import {VStack, Text} from '@chakra-ui/react'

export type Props = {title: string, description?: string}

const Label: React.FC<Props> = ({title, description, children}) => {
    return (
        <VStack alignItems="flex-start" spacing=".15rem">
            <Text color="text.600" fontSize="1rem" fontWeight={600} letterSpacing="wide" _after={{content: '":"'}}>{title}</Text>
            {description && <Text color="text.700" fontSize="0.95rem">{description}</Text>}
            {children && children}
        </VStack>
    )
}

export default Label