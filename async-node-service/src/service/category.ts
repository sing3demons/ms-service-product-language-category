import dayjs from 'dayjs'
import { Category } from '../models/category.js'
import { insertOneCategory, findOneCategory, findCategoryId, updateProduct } from '../repository/category.js'
import logger from '../utils/logger.js'
import NanoIdService from '../utils/nanoid.js'

async function createCategory(req: Category) {
  const nano = new NanoIdService()
  try {
    req.id = req.id || nano.randomNanoId()
    req.lastUpdate = new Date().toISOString()
    if (req.validFor && Object.keys(req.validFor).length !== 0) {
      req.validFor.startDateTime = dayjs(req.validFor?.startDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
      req.validFor.endDateTime = dayjs(req.validFor?.endDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
    }
    const doc: Category = {
      id: req.id,
      '@type': 'Category',
      lastUpdate: req.lastUpdate,
      name: req.name,
      products: req.products || [],
      version: req.version,
      validFor: req.validFor,
    }

    const result = await insertOneCategory(doc)
    const category = await findCategoryId(result.insertedId)
    return category
  } catch (e) {
    if (e instanceof Error) {
      logger.error(e.message)
      throw new Error(e.message)
    }
  }
}

async function updateCategory(req: Category) {
  try {
    if (!req.id) {
      throw new Error('ID is required')
    }

    const result = await findOneCategory(req.id)
    if (!!result) {
      if (req.products) {
        if (Array.isArray(result.products)) {
          for (let i = 0; i < req.products.length; i++) {
            const product = req.products[i]
            if (product?.id) {
              result.products.push({ id: product.id, name: product?.name })
            }
          }

          const productIds = new Set(result.products.map((item) => item.id))
          const products = req.products.filter((item) => !productIds.has(item.id))
          console.log(products)

          const data = await updateProduct(req.id, products)
          return data
        }
      }
    }
  } catch (e) {
    if (e instanceof Error) {
      logger.error(e.message)
      throw new Error(e.message)
    }
  }
}

export { createCategory, updateCategory }
