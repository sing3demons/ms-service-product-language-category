import { ObjectId, WithId } from 'mongodb'
import { getCollection } from '../db.js'

import logger from '../utils/logger.js'
import { ProductLanguage } from '../models/productLanguage.js'

const col = getCollection<ProductLanguage>('productLanguage')

async function insertOneProductLanguage(req: ProductLanguage) {
  try {
    logger.info('insertOneProductLanguage', req)
    const result = await col.insertOne(req)
    return result
  } catch (e) {
    throw e
  }
}

async function findProductLanguageId(id: ObjectId | string) {
  let result: WithId<ProductLanguage> | null = null
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

export { insertOneProductLanguage, findProductLanguageId }
