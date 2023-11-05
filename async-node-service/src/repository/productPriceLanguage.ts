import { ObjectId, WithId } from 'mongodb'
import { getCollection } from '../db.js'

import logger from '../utils/logger.js'
import { IProductPriceLanguageDTO } from '../models/productPriceLanguage.js'

const col = getCollection<IProductPriceLanguageDTO>('productPriceLanguage')

async function insertOneProductPriceLanguage(req: IProductPriceLanguageDTO) {
  try {
    logger.info('insertOneProductPrice', req)
    const result = await col.insertOne(req)
    return result
  } catch (e) {
    throw e
  }
}

async function findProductPriceLanguageId(id: ObjectId | string) {
  let result: WithId<IProductPriceLanguageDTO> | null = null
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

export { insertOneProductPriceLanguage, findProductPriceLanguageId }
