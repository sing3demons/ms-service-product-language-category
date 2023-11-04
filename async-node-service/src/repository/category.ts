import { ObjectId, WithId } from 'mongodb'
import { getCollection } from '../db.js'
import { Category } from '../models/category.js'
import logger from '../utils/logger.js'
import { Product } from '../models/product.js'

const col = getCollection<Category>('category')

async function insertOneCategory(req: Category) {
  try {
    const result = await col.insertOne(req)
    return result
  } catch (e) {
    throw e
  }
}

async function findCategoryId(_id: ObjectId) {
  try {
    const result = await col.findOne({ _id })
    return result
  } catch (e) {
    throw e
  }
}

async function findOneCategory(id: string) {
  try {
    const result: WithId<Category> | null = await col.findOne({ id: id })
    return result as Category
  } catch (e) {
    throw e
  }
}

async function updateCategory(id: string, req: Category) {
  try {
    const result = await col.updateOne({ id }, { $set: req })
    return result
  } catch (e) {
    throw e
  }
}

async function updateProduct(id: string, products: Product[]) {
  try {
    console.log('=================')
    console.log(products)
    const result = await col.findOneAndUpdate({ id }, { $set: { products } }, { upsert: true, returnDocument: 'after' })
    console.log('=======  result  ==========')
    console.log(JSON.stringify(result))
    return result
  } catch (e) {
    throw e
  }
}

export { insertOneCategory, findOneCategory, findCategoryId, updateProduct }
