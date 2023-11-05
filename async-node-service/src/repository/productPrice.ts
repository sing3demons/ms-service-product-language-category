import { ObjectId, WithId } from 'mongodb'
import { getCollection } from '../db.js'

import logger from '../utils/logger.js'
import { IProductPriceDTO } from '../models/productPrice.js'

const col = getCollection<IProductPriceDTO>('productPrice')

async function insertOneProductPrice(req: IProductPriceDTO) {
  try {
    logger.info('insertOneProductPrice', req)
    const result = await col.insertOne(req)
    return result
  } catch (error) {
    throw error
  }
}

async function findProductPriceId(id: ObjectId | string) {
  let result: WithId<IProductPriceDTO> | null = null
  try {
    if (id instanceof ObjectId) {
      result = await col.findOne({ _id: id })
    } else {
      result = await col.findOne({ id })
    }
    return result
  } catch (error) {
    throw error
  }
}

export { insertOneProductPrice, findProductPriceId }
