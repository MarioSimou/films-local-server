import { S3Client, DeleteObjectCommand, CreateMultipartUploadCommand, UploadPartCommand, UploadPartRequest, CompleteMultipartUploadCommand} from '@aws-sdk/client-s3'
import type { DeleteObjectCommandInput, CreateMultipartUploadCommandInput} from '@aws-sdk/client-s3'
import type { HTTPResponse } from '@libs/types'

const bucketName = process.env.AWS_S3_BUCKET_NAME
const client = new S3Client({
    endpoint: process.env.AWS_ENDPOINT,
    forcePathStyle: true
})

type Body = UploadPartRequest["Body"] | string | Uint8Array | Buffer;
export const uploadOne = async (key: CreateMultipartUploadCommandInput['Key'], body: Body, options: Partial<Omit<CreateMultipartUploadCommandInput, 'Key'>> = {}): Promise<HTTPResponse<string>> => {
    try {
        const {UploadId: uploadId } = await client.send(new CreateMultipartUploadCommand({
            Bucket: bucketName,
            Key: key,
            StorageClass: "STANDARD",
            ...options,
        }))

        const partNumber = 1
        const {ETag: etag} = await client.send(new UploadPartCommand({
            Bucket: bucketName,
            Key: key,
            Body: body,
            UploadId: uploadId,
            PartNumber: partNumber,
        }))

        const {Location} = await client.send(new CompleteMultipartUploadCommand({
            Bucket: bucketName,
            UploadId: uploadId, 
            Key: key,
            MultipartUpload: {
                Parts: [
                    {ETag: etag, PartNumber: partNumber }
                ]
            }
        }))
        return [undefined, Location]    
    }catch(e){
        return [e]
    }
}

export const deleteOne = async (key: DeleteObjectCommandInput['Key'] ): Promise<[undefined | Error]> => {
    try {
        const command = new DeleteObjectCommand({
            Bucket: bucketName,
            Key: key
        })
        await client.send(command)
        return [undefined]
    }catch(e){
        return [e]
    }
}