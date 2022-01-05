import React from 'react'
import {Flex, Input as ChakraInput, Text, Icon, Collapse, Heading, VStack, Img, useToast} from '@chakra-ui/react'
import {TextareaField, InputField} from '@components/shared/Field'
import { useFormFields, useSong } from '@hooks'
import {FaCloudUploadAlt as UploadIcon} from 'react-icons/fa'
import {httpClient} from '@utils'
import type { HTTPResponse, Song } from '@types'
import type { AxiosResponse } from 'axios'
import { useRouter } from 'next/router'
import {ChakraButton, Button} from '@components/shared/Button'

const Upload = () => {
    const {isLoading, postSong} = useSong()
    const router = useRouter()
    const toast = useToast({
        status: "error",
        title: "Upload",
        position: 'bottom-right',
        isClosable: true,
    })
    const {fields, handleOnBlur, handleOnChange} = useFormFields<string>({
        name: {
            touched: false,
            error: '',
            value: '',
        },
        description: {
            touched: false,
            error: '',
            value: '',
        },
        image: {
            touched: false,
            error: '',
            value: '',
        }
    })
    const imageRef = React.useRef<HTMLInputElement>()
    const locationRef = React.useRef<HTMLInputElement>()
    const {fields: fileFields, setValue: setFileValue} = useFormFields<File | undefined>({
        image: {
            touched: false,
            error: '',
            value: imageRef.current?.files?.[0],
        },
        location: {
            touched: false,
            error: '',
            value: locationRef.current?.files?.[0],
        },
    })
    const onClickUploadFile = React.useCallback((ref: React.MutableRefObject<HTMLInputElement | undefined>) => {
        return () => ref.current?.click()
    }, [])

    const onChangeFile = React.useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const {id} = e.currentTarget 
        const file = e.currentTarget.files?.[0]

        if(!file){
            return
        }

        setFileValue(id, file, '')
    }, [setFileValue])

    const onUploadSong = React.useCallback(async () => {
        if(!fields.name.value || !fields.description.value || !fileFields.image.value || !fileFields.location.value) {
            return toast({description: 'Please make sure that you set a name, image and a song video.'})
        }

        const [e] = await postSong({
            name: fields.name.value,
            description: fields.description.value,
            image: fileFields.image.value,
            location: fileFields.location.value,
        })

        if(e){
            return toast({description: e.message})
        }

        toast({description: "The song has been succesfully uploaded!", status: "success"})
        router.push('/')
    }, [fields.name.value, fields.description.value, fileFields.image.value, fileFields.location.value, toast, router, postSong])

    return (
        <Flex flexDirection="column" alignItems="center" minH="calc(100vh - 84px)" bg="gray.100">
            <VStack spacing="2rem" alignItems="flex-start" maxWidth="1000px" w="100%" p="4rem 2rem">
                <VStack w="100%">
                    <Heading fontSize="1.75rem">Upload a song</Heading>
                    <Text color="text.600">This page is used to upload a song to our platform</Text>
                </VStack>
                <VStack w="100%">
                    <VStack as="form" spacing="3rem" w="100%">
                        <InputField id="name" label="Name" placeholder="Your song name" value={fields.name.value} error={fields.name.error} touched={fields.name.touched} onChange={handleOnChange} onBlur={handleOnBlur}/>
                        <TextareaField id="description" label="Description" placeholder="A few words about the song" value={fields.description.value} error={fields.description.error} touched={fields.description.touched} onChange={handleOnChange} onBlur={handleOnBlur}/>
                        <VStack bg="white" w="100%" p="4rem 2rem" borderRadius="1rem" spacing="0.75rem">
                            <Icon as={UploadIcon} w={12} h={12}/>
                            <Text fontSize="0.9rem" color="text.700">Select an image to upload</Text>
                            <ChakraButton variant="outline" colorScheme="primary" color="primary.900" onClick={onClickUploadFile(imageRef)}>Upload an image</ChakraButton>
                            <ChakraInput type="file" id="image" name="image" ref={imageRef as any} display="none" accept="image/*" onChange={onChangeFile}/>
                            <Collapse in={Boolean(fileFields.image.value)}>
                                {fileFields.image.value && <Img src={URL.createObjectURL(fileFields.image.value)} alt="song-image" w="300px" h="200px"/>}
                            </Collapse>
                        </VStack>
                        <VStack bg="white" w="100%" p="4rem 2rem" borderRadius="1rem" spacing="0.75rem">
                            <Icon as={UploadIcon} w={12} h={12}/>
                            <Text fontSize="0.9rem" color="text.700">Select a song to upload</Text>
                            <ChakraButton variant="outline" colorScheme="primary" color="primary.900"  onClick={onClickUploadFile(locationRef)}>Upload a song</ChakraButton>
                            <ChakraInput type="file" id="location" ref={locationRef as any} display="none" accept="video/*" onChange={onChangeFile}/>
                            <Collapse in={Boolean(fileFields.location.value)}>
                                {fileFields.location.value && <video width="360px" height="240px" controls>
                                    <source src={URL.createObjectURL(fileFields.location.value)} type={fileFields.location.value?.type}/>
                                </video>}
                            </Collapse>
                        </VStack>
                        <Button isLoading={isLoading} onClick={onUploadSong} alignSelf="flex-start">Upload</Button>
                    </VStack>
                </VStack>
            </VStack>   
        </Flex>
    )
}

export default Upload