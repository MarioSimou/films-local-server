import { formatJSONResponse } from '@libs/api-gateway'

const main = async () => {
    return formatJSONResponse({message: 'upload song'})
}

export default main