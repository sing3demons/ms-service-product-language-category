import { IProductPriceDTO } from '../models/productPrice.js'
import { findProductPriceId, insertOneProductPrice } from '../repository/productPrice.js'

async function createProductPrice(req: IProductPriceDTO) {
  try {
    const doc: IProductPriceDTO = {
      '@type': 'productPrice',
      id: req.id,
      lastUpdate: req.lastUpdate,
      name: req.name,
      version: req.version,
      price: req.price,
      validFor: req.validFor,
      lifecycleStatus: req.lifecycleStatus,
      popRelationship: req.popRelationship,
      TORO_supportingLanguage: req.TORO_supportingLanguage,
    }

    if (doc.validFor && Object.keys(doc.validFor).length === 0) {
      delete doc.validFor
    }

    const result = await insertOneProductPrice(doc)
    return await findProductPriceId(result.insertedId)
  } catch (error) {
    throw error
  }
}

export { createProductPrice }
