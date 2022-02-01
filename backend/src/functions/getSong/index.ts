import { formatJSONResponse } from '@libs/api-gateway'
const main = async () => {
    return formatJSONResponse({message: 'get song'})
}

export default main