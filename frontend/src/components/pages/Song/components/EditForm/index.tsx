import React from 'react'
import { Modal, ModalOverlay, ModalHeader, Img, ModalCloseButton, ModalBody, ModalFooter, ModalContent, VStack, useToast} from '@chakra-ui/react'
import type { Song  } from '@types'
import { useFormFields, useSong } from '@hooks'
import {InputField, TextareaField} from '@components/shared/Field'
import { Button, ChakraButton } from '@components/shared/Button'
import FileInput from '@components/shared/FileInput'

export type EditFormProps = {
    isOpen: boolean
    onClose: () => void
    song: Omit<Song, 'href' | 'image'> & {href: File, image: File},
}

const EditForm = ({isOpen, onClose, song}: EditFormProps) => {
    const {isLoading, putSong} = useSong()
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

    const { fields: fileFields, setValue: setFileValue } = useFormFields({
        image: {
            value: song.image,
            touched: false,
            error: '',
        },
        href: {
            value: song.href,
            touched: false,
            error: ''
        }
    })

    const onChangeFile = React.useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const {id} = e.target 
        const file = e.currentTarget.files?.[0]

        if(!file){
            return
        }

        setFileValue(id, file, '')
    }, [setFileValue])
    
    const onUpdateSong = async () => {
        const isTouched = fields.name.touched || fields.description.touched || fileFields.image.touched || fileFields.href.touched
 
        if(!isTouched){
            return toast({description: "We didn't notice any change on the song structure. Please provide some changes."})
        }
        const [e] = await putSong(song.guid, {
            name: fields.name.value,
            description: fields.description.value,
            image: fileFields.image.value,
            href: fileFields.href.value,
        })

        if(e){
            return toast({description: e.message})
        }

        toast({description: "The song has been succesfully updated", status: "success"})
        return onClose()
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
                        {song.image && <FileInput bg="gray.50" id="image" subtitle="Select an image to upload" btnText="Upload an image" value={fileFields.image.value} onChangeFile={onChangeFile}>
                                <Img w="360px" h="240px" src={URL.createObjectURL(fileFields.image.value)} type={fileFields.image.value?.type}/>
                            </FileInput>}
                        {song.href && <FileInput bg="gray.50" accept="audio/*" id="href" subtitle="Select a song to upload" btnText="Upload a song" value={fileFields.href.value} onChangeFile={onChangeFile}>
                            <audio controls>
                                <source src={URL.createObjectURL(fileFields.href.value)} type={fileFields.href.value?.type}/>
                            </audio>
                        </FileInput>}
                    </VStack>
                </ModalBody>
                <ModalFooter>
                    <VStack w="100%">
                        <Button isLoading={isLoading} onClick={onUpdateSong}>Update</Button>
                        <ChakraButton isFullWidth variant="outline" colorScheme="primary" color="primary.900" onClick={onClose}>Close</ChakraButton>
                    </VStack>
                </ModalFooter>
            </ModalContent>
        </Modal>
    )
}

export default EditForm