import React from 'react'
import {VStack, Icon, Input as ChakraInput, Button as ChakraButton, Collapse, Text, StackProps} from '@chakra-ui/react'
import {FaCloudUploadAlt as UploadIcon} from 'react-icons/fa'

export type FileInputProps = {
    id: string
    value: File | undefined
    subtitle?: string
    btnText: string
    onChangeFile: (e: React.ChangeEvent<HTMLInputElement>) => void
    accept?: string
} & StackProps

const FileInput: React.FC<FileInputProps> = ({children, value, btnText, subtitle, onChangeFile, id, accept = "image/*", ...other}) => {
    const ref = React.useRef<HTMLInputElement>()
    const onClickButton = React.useCallback(() => {
        return ref.current?.click()
    }, [ref])

    return (
        <VStack bg="white" w="100%" p="2rem" borderRadius="1rem" spacing="0.75rem" {...other}>
            <Icon as={UploadIcon} w={12} h={12}/>
            {subtitle && <Text fontSize="0.9rem" color="text.700">{subtitle}</Text>}
            <ChakraButton variant="outline" colorScheme="primary" color="primary.900"  onClick={onClickButton}>{btnText}</ChakraButton>
            <ChakraInput type="file" id={id} name={id} ref={ref as any} display="none" accept={accept} onChange={onChangeFile}/>
            {children && <Collapse in={Boolean(value)} style={{marginTop: "2rem"}}>
                {children}
            </Collapse>}
        </VStack>
    )
}

export default FileInput