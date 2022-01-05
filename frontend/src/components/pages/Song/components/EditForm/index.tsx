import React from 'react'
import { Modal, ModalOverlay, ModalHeader, ModalCloseButton, ModalBody, ModalFooter, ModalContent, VStack, useToast} from '@chakra-ui/react'
import type { Song  } from '@types'
import { useFormFields } from '@hooks'
import {InputField, TextareaField} from '@components/shared/Field'
import { Button, ChakraButton } from '@components/shared/Button'
import { httpClient } from '@utils'

export type EditFormProps = {
    isOpen: boolean
    onClose: () => void
    song: Song
}

const fetchFile = async (name: string, href: string): Promise<[Error] | [undefined, File]> => {
    try {
        const res = await fetch(href)
        const blob = await res.blob()
        const file = new File([blob], name)
        return [undefined, file]
    
    } catch(e){
        return [e.message]
    }
}

const EditForm = ({isOpen, onClose, song}: EditFormProps) => {
    const toast = useToast({
        isClosable: true,
        title: 'Update Song',
        status: 'error',
        position: 'bottom-right',
    })
    const {fields, handleOnChange, handleOnBlur} = useFormFields({
        name: {
            value: song.name,
            touched: false,
            error: '',
        },
        description: {
            value: song.description ?? '',
            touched: false,
            error: '',
        },
    })
    const imageRef = React.useRef<File>()
    const locationRef = React.useRef<File>()

    const { fields: fileFields, setValue } = useFormFields({
        image: {
            value: imageRef.current,
            touched: false,
            error: '',
        },
        location: {
            value: locationRef.current,
            touched: false,
            error: ''
        }
    })
    
    // React.useEffect(() => {
    //     const fetchFiles = async () => {
    //         const [imageError, imageFile] = await fetchFile("image",song.image)
    //         if(imageError){
    //             return toast({description: imageError.message})
    //         }

    //         const [locationError, locationFile] = await fetchFile("location",song.location)
    //         if(locationError){
    //             return toast({description: locationError.message})
    //         }

    //         setValue('image', imageFile, '')
    //         setValue('location', locationFile, '')
    //     }
    //     fetchFiles()
    // }, [])

    const onUpdateSong = () => {
        const isTouched = fields.name.touched || fields.description.touched

        if(!isTouched){
            return toast({description: "We didn't notice any change on the song structure. Please provide some changes."})
        }

        console.log('update song')
    }

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="full">
            <ModalOverlay/>
            <ModalContent borderRadius="0">
                <ModalHeader>Header</ModalHeader>
                <ModalCloseButton/>
                <ModalBody>
                    <VStack spacing="1rem">
                        <InputField id="name" label="Name" onChange={handleOnChange} onBlur={handleOnBlur} {...fields.name}/>
                        <TextareaField id="description" label="Description" onChange={handleOnChange} onBlur={handleOnBlur} {...fields.description}/>
                    </VStack>
                </ModalBody>
                <ModalFooter>
                    <VStack w="100%">
                        <Button onClick={onUpdateSong}>Update</Button>
                        <ChakraButton isFullWidth variant="outline" colorScheme="primary" color="primary.900" onClick={onClose}>Close</ChakraButton>
                    </VStack>
                </ModalFooter>
            </ModalContent>
        </Modal>
    )
}

export default EditForm