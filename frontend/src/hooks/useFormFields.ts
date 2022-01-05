import React from 'react'

type FieldOptions<T> = {
    value: T
    error: string
    touched: boolean
}

export const useFormFields = <T>(initialFields: Record<string, FieldOptions<T>>) => {
    const [fields, setFields] = React.useState(initialFields)

    const setTouched = React.useCallback((id: string, touched = true) => {
        return setFields(fields => ({
            ...fields,
            [id]: {
                ...fields[id],
                touched
            }
        }))
    }, [setFields])

    const setValue = React.useCallback((id: string, value: T, error: string) => {
        return setFields(fields => ({
            ...fields,
            [id]: {
                ...fields[id],
                touched: true,
                error,
                value,
            }
        }))
    }, [])

    const setError = React.useCallback((id: string, error: string) => {
        return setFields(fields => ({
            ...fields,
            [id]: {
                ...fields[id],
                touched: true,
                error,
            }
        }))
    }, [])

    const handleOnChange = React.useCallback(<I extends HTMLInputElement | HTMLTextAreaElement>(e: React.ChangeEvent<I>) => {
        const value: string = e.currentTarget.value
        const id = e.currentTarget.id

        if(!value){
            return setValue(id, value as any, 'This field is required.')
        }
        return setValue(id, value as any, '')
    }, [setValue])

    const handleOnBlur = React.useCallback(<T extends HTMLInputElement | HTMLTextAreaElement>(e: React.ChangeEvent<T>) => {
        return handleOnChange(e)
    }, [handleOnChange])


    return {
        fields,
        setError,
        setValue, 
        setTouched,
        handleOnBlur,
        handleOnChange,
    }
}