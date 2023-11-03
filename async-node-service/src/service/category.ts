import dayjs from 'dayjs'
import { Category } from '../models/category.js'
import { insertOneCategory, findOneCategory, findCategoryId, updateProduct } from '../repository/category.js'
import logger from '../utils/logger.js'
import NanoIdService from '../utils/nanoid.js'
import { Product } from '../models/product.js'
import { findProductByIds } from '../repository/product.js'

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
      lifecycleStatus: req.lifecycleStatus,
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
      logger.info(JSON.stringify(result))
      if (req.products) {
        if (Array.isArray(result.products) && result.products.length !== 0) {
          const productIds = await findProductByIds(req.products.map((item) => item.id))
          console.log('====== productIds ===========')
          console.log(JSON.stringify(productIds))
          const productMap = new Map()
          if (productIds.length !== 0) {
            for (let i = 0; i < productIds.length; i++) {
              const product = productIds[i]
              if (product?.id) {
                productMap.set(product.id, product)
              }
            }

            let products: Product[] = []
            for (let i = 0; i < req.products.length; i++) {
              const product = req.products[i]
              if (product?.id) {
                if (productMap.has(product.id)) {
                  logger.info(`has product id ${product.id} ${product.name}`)
                  productMap.delete(product.id)
                }
                const update = { id: product.id, name: product?.name }
                productMap.set(product.id, update)
                logger.info('for update', update.id, update.name)
                products.push(update)
              }
            }

            console.log('====== category.products ===========')
            const data = await updateProduct(req.id, [...productMap.values()])
            return data
          }
        } else {
          console.log('====== new category.products ===========')
          console.log(req.products)
          const products: Product[] = req.products.map((item) => {
            return { id: item.id, name: item.name }
          })
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
