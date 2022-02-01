import { formatJSONResponse } from '@libs/api-gateway'
const main = async () => {
    return formatJSONResponse({message: 'delete song'})
}

export default main