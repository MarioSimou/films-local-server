import { DynamoDBClient, PutItemCommand, GetItemCommand, GetItemOutput, ScanCommand, DeleteItemCommand } from '@aws-sdk/client-dynamodb'
import type { Song, HTTPResponse } from '@libs/types'
import { v4 as uuidv4 } from 'uuid'
import { ErrNotFound } from '@libs/api-gateway'

const dbClient = new DynamoDBClient({endpoint: process.env.AWS_ENDPOINT})
const songsTableName = process.env.AWS_DYNAMODB_SONGS_TABLE_NAME

const toSong = (item: GetItemOutput['Item']): Song => ({
    guid: item?.guid?.S,
    name: item?.name?.S,
    description: item?.description?.S,
    image: item?.image?.S,
    href: item?.href?.S,
    createdAt: item?.createdAt?.S,
})

const postPayloadToItem = ({name, description, guid}: Pick<Song, 'guid' | 'name' | 'description'>) => ({
    guid: {
        'S': guid,
    },
    name: {
        'S': name,
    },
    description: {
        'S': description
    },
    createdAt: {
        'S': new Date().toISOString()
    }
})

const putPayloadToItem = (payload: Song) => {
    return Object.entries(payload).reduce((item, [key, value]) => {
        if(!value){
            return item
        }
        return {...item, [key]: {S: value}}
    },{})
}

const toGUIDItem = (guid: string) => ({
    'guid': {
        'S': guid
    }
})

export const postSong = async (payload :Pick<Song, 'name' | 'description'>): Promise<HTTPResponse<Song>> => {
    try {
        const guid = uuidv4()
        const command = new PutItemCommand({
            TableName: songsTableName,
            Item: postPayloadToItem({...payload, guid}),
        })
    
        await dbClient.send(command)

        const [e, song] = await findOneSong(guid)
    
        if(e){
            return [e]
        }

        return [undefined, song]
    }catch(e){
        return [e as Error]
    }
}

export const findOneSong = async (guid: string): Promise<HTTPResponse<Song>> => {
   try {
    const command = new GetItemCommand({
        TableName: songsTableName,
        Key: {
            guid: {
                S: guid,
            }
        }
    })

    const res = await dbClient.send(command)

    if(!res.Item){
        return [ErrNotFound]
    }
    
    return [undefined, toSong(res.Item)]
   }catch(e){
       return [e as Error]
   }
}

export const findAllSongs = async (limit: number = 20): Promise<HTTPResponse<Song[]>> => {
    try {
        const command = new ScanCommand({TableName: songsTableName, Limit: limit})
        const {Items} = await dbClient.send(command)
        const songs = Items.map(toSong)
        return [undefined, songs]
    }catch(e){
        return [e as Error]
    }
}

export const putOneSong = async (guid: string, payload: Song): Promise<HTTPResponse<Song>> => {
    try {
        const command = new PutItemCommand({
            TableName: songsTableName,
            Item: putPayloadToItem(payload),
        })
        await dbClient.send(command)
        const [e, song] = await findOneSong(guid)
        if(e){
            return [e]
        }
        return [undefined, song]
    }catch(e){
        return [e as Error]
    }
}

export const deleteOne = async (key: string): Promise<[Error|undefined]> => {
    try {
        const command = new DeleteItemCommand({
            TableName: songsTableName,
            Key: toGUIDItem(key),
        })
        await dbClient.send(command)
        return [undefined]
    }catch(e){
        return [e as Error]
    }
}