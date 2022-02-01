import { formatJSONResponse } from '@libs/api-gateway'
const main = async () => {
    return formatJSONResponse({message: 'post song'})
}

export default main