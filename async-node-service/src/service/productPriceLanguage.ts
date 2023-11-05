import { IProductPriceLanguageDTO } from '../models/productPriceLanguage.js'
import { findProductPriceLanguageId, insertOneProductPriceLanguage } from '../repository/productPriceLanguage.js'

async function createProductPriceLanguage(req: IProductPriceLanguageDTO) {
  try {
    const doc: IProductPriceLanguageDTO = {
      '@type': 'productPriceLanguage',
      id: req.id,
      languageCode: req.languageCode,
      lastUpdate: req.lastUpdate,
      name: req.name,
      version: req.version,
      price: req.price,
      validFor: req.validFor,
    }

    if (doc.validFor && Object.keys(doc.validFor).length === 0) {
      delete doc.validFor
    }

    const result = await insertOneProductPriceLanguage(doc)
    return await findProductPriceLanguageId(result.insertedId)
  } catch (error) {
    throw error
  }
}

export { createProductPriceLanguage }
