import {
  S3Client,
  DeleteObjectCommand,
  CreateMultipartUploadCommand,
  UploadPartCommand,
  CompleteMultipartUploadCommand,
} from '@aws-sdk/client-s3'
import type {Result} from '@libs/types'

const AWS_ENDPOINT = process.env.AWS_ENDPOINT
const AWS_S3_ASSETS_BUCKET_NAME = process.env.AWS_S3_ASSETS_BUCKET_NAME

if (!AWS_S3_ASSETS_BUCKET_NAME) {
  throw new Error('error: s3 configuration is invalid')
}

const getS3Client = () => {
  if (!AWS_ENDPOINT) {
    return new S3Client({forcePathStyle: true})
  }
  return new S3Client({
    endpoint: AWS_ENDPOINT,
    forcePathStyle: true,
  })
}

const deleteObject =
  (client: S3Client) =>
  async (key: string): Promise<[Error | undefined]> => {
    try {
      const command = new DeleteObjectCommand({
        Bucket: AWS_S3_ASSETS_BUCKET_NAME,
        Key: key,
      })
      await client.send(command)
      return [undefined]
    } catch (e) {
      return [e]
    }
  }

const uploadObject =
  (client: S3Client) =>
  async (key: string, contentType: string, body: string): Promise<Result<string>> => {
    try {
      const createMultipartUploadCommand = new CreateMultipartUploadCommand({
        Bucket: AWS_S3_ASSETS_BUCKET_NAME,
        StorageClass: 'STANDARD',
        ContentType: contentType,
        Key: key,
        ACL: 'public-read',
      })
      const {UploadId: uploadId} = await client.send(createMultipartUploadCommand)

      const uploadPartCommand = new UploadPartCommand({
        Bucket: AWS_S3_ASSETS_BUCKET_NAME,
        Body: Buffer.from(body, 'base64'),
        UploadId: uploadId,
        PartNumber: 1,
        Key: key,
      })

      const {ETag: etag} = await client.send(uploadPartCommand)

      const completeMultipartUploadCommand = new CompleteMultipartUploadCommand({
        UploadId: uploadId,
        Bucket: AWS_S3_ASSETS_BUCKET_NAME,
        Key: key,
        MultipartUpload: {
          Parts: [{ETag: etag, PartNumber: 1}],
        },
      })
      const {Location: location} = await client.send(completeMultipartUploadCommand)
      return [undefined, location]
    } catch (e) {
      return [e]
    }
  }
type S3ClientFactoryResult = {
  uploadObject: ReturnType<typeof uploadObject>
  deleteObject: ReturnType<typeof deleteObject>
}

export const newS3ClientFactory = (): S3ClientFactoryResult => {
  const client = getS3Client()
  return {
    uploadObject: uploadObject(client),
    deleteObject: deleteObject(client),
  }
}
