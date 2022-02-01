import { formatJSONResponse } from '@libs/api-gateway'
const main = async () => {
    return formatJSONResponse({message: 'put song'})
}

export default main