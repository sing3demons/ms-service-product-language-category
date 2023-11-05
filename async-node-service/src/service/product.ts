import { Category } from '../models/category.js'
import { Product } from '../models/product.js'
import { ProductLanguage } from '../models/productLanguage.js'
import { findProductId, insertOneProduct } from '../repository/product.js'
import dayjs from 'dayjs'

async function createProduct(req: Product) {
  try {
    const category: Category[] = []
    if (Array.isArray(req.category)) {
      for (let i = 0; i < req.category.length; i++) {
        const item = req.category[i]
        if (item?.id) {
          category.push({
            id: item.id,
            name: item?.name,
            '@type': 'Category',
          })
        }
      }
    }

    const supportingLanguage: ProductLanguage[] = []
    if (Array.isArray(req.supportingLanguage)) {
      for (let i = 0; i < req.supportingLanguage.length; i++) {
        const item = req.supportingLanguage[i]
        if (item?.id) {
          supportingLanguage.push({
            id: item.id,
            name: item?.name || undefined,
            '@type': 'ProductLanguage',
            languageCode: item.languageCode,
            referredType: item['referredType'],
          })
        }
      }
    }

    if (req.validFor && Object.keys(req.validFor).length !== 0) {
      req.validFor.startDateTime = dayjs(req.validFor?.startDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
      req.validFor.endDateTime = dayjs(req.validFor?.endDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
    }
    const doc: Product = {
      id: req.id,
      '@type': 'Product',
      category: category,
      description: req.description,
      lastUpdate: req.lastUpdate,
      lifecycleStatus: req.lifecycleStatus,
      name: req.name,
      price: req.price,
      supportingLanguage: supportingLanguage,
      validFor: req.validFor,
      version: req.version,
    }

    if (doc.validFor && Object.keys(doc.validFor).length === 0) {
      delete doc.validFor
    }
    const result = await insertOneProduct(doc)
    const data = await findProductId(result.insertedId)
    return data
  } catch (e) {
    throw e
  }
}

export { createProduct }
