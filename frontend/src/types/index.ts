export type ITimestamp = {
    createdAt: string
    updatedAt: string
}
export type IGUID = {
    guid: string
}

export type Song = {
    name: string
    description?: string
    image: string
    href: string
} & ITimestamp & IGUID

export type PostSong = {
    name: string
    description: string
    image: File
    href: File
}
export type PutSong = PostSong

export type HTTPResponse<T> = {
    status: number
    success: boolean
    message?: string
    data?: T
}