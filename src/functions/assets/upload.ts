import {SNSEvent} from 'aws-lambda'
import aws, {PublishUploadAssetMessage} from '@libs/aws'

export const handler = async (event: SNSEvent): Promise<void> => {
  const [record] = event.Records
  const {Message: message} = record.Sns
  const {asset: body, key, contentType, song, field}: PublishUploadAssetMessage = JSON.parse(message)

  const [uploadObjectError, location] = await aws.s3.uploadObject(key, contentType, body)
  if (uploadObjectError) {
    return console.error(uploadObjectError.message)
  }

  const [putSongError] = await aws.dynamoDB.putSong({
    ...song,
    [field]: location,
    [`${field}Status`]: 'completed',
  })

  if (putSongError) {
    return console.error(uploadObjectError.message)
  }

  return undefined
}
