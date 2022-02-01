import { formatJSONResponse } from '@libs/api-gateway'
const main = async () => {
    return formatJSONResponse({message: 'put image song'})
}

export default main