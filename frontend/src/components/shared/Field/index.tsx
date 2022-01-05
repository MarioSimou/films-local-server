import React from 'react'
import {Input as ChakraInput, FormControl,FormLabel, FormErrorMessage, Collapse, Textarea as ChakraTextarea, TextareaProps as ChakraTextareaProps, InputProps as ChakraInputProps} from '@chakra-ui/react'

export type Props<P = ChakraInputProps | ChakraTextareaProps> = {
    id: string
    label?: string
    value: string
    error: string
    touched: boolean
    required?: boolean
    placeholder?: string
    InputProps?: P
    children: React.ReactElement<P>
    onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void
    onBlur?: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void
}

export const Field = function<P = ChakraInputProps | ChakraTextareaProps>({
    id,
    label = '',
    error,
    touched,
    required = true,
    children,
}: Omit<Props<P>, 'placeholder' | 'onChange' | 'onBlur' | 'value' | 'InputProps'>){
    const isInvalid = Boolean(touched && error)
    return (
        <FormControl isRequired={required} isInvalid={isInvalid} id={id}>
            {label && <FormLabel mb="0rem" color="text.700" letterSpacing="wide">{label}</FormLabel>}
            {children}
            <Collapse in={isInvalid}>
                <FormErrorMessage>{error}</FormErrorMessage>
            </Collapse>
        </FormControl>
    )
}

export const InputField = function({
    id,
    value,
    placeholder,
    onBlur,
    onChange,
    InputProps,
    ...rest
}: Omit<Props, 'children'>){
    return (
        <Field id={id} {...rest}>
            <ChakraInput variant="flushed" name={id} value={value} onChange={onChange} onBlur={onBlur} placeholder={placeholder} {...InputProps as any}/>
        </Field>
    )
}

export const TextareaField = function({
    id,
    value,
    placeholder,
    onBlur,
    onChange,
    InputProps,
    ...rest
}: Omit<Props, 'children'>){
    return (
        <Field id={id} {...rest}>
            <ChakraTextarea variant="flushed" name={id} value={value} onChange={onChange} onBlur={onBlur} placeholder={placeholder} {...InputProps as any}/>
        </Field>
    )
}

export default Field