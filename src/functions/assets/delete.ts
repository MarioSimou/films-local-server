import {SNSEvent} from 'aws-lambda'
import aws, {PublishDeleteAssetMessage} from '@libs/aws'
import {extname} from 'path'

export const handler = async (event: SNSEvent): Promise<void> => {
  const [record] = event.Records
  const {Message: message} = record.Sns
  const {song, field}: PublishDeleteAssetMessage = JSON.parse(message)

  const key = `${song.guid}/${field}${extname(song[field])}`
  const [deleteObjectError] = await aws.s3.deleteObject(key)
  if (deleteObjectError) {
    return console.error(deleteObjectError.message)
  }

  return undefined
}
