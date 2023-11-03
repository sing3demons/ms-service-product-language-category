import { ObjectId, WithId } from 'mongodb'
import { getCollection } from '../db.js'
import { Product } from '../models/product.js'

const col = getCollection<Product>('product')

async function insertOneProduct(req: Product) {
  try {
    const result = await col.insertOne(req)
    return result
  } catch (e) {
    throw e
  }
}

async function findProductId(_id: ObjectId | string) {
  let result: WithId<Product> | null = null
  try {
    if (_id instanceof ObjectId) {
      result = await col.findOne({ _id })
    } else {
      result = await col.findOne({ id: _id })
    }
    return result
  } catch (error) {
    throw error
  }
}

async function findProductByIds(ids: string[]) {
  console.log(`findProductByIds ${ids}`)
  try {
    const results = await col.find({ id: { $in: ids } }).toArray()
    return results
  } catch (error) {
    throw error
  }
}

export { insertOneProduct, findProductId, findProductByIds }
