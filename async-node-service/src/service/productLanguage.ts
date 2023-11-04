import dayjs from 'dayjs'
import { Attachment, ProductLanguage } from '../models/productLanguage.js'
import { findProductLanguageId, insertOneProductLanguage } from '../repository/productLanguage.js'
import NanoIdService from '../utils/nanoid.js'

async function createProductLanguage(req: ProductLanguage) {
  const nano = new NanoIdService()
  let attachments: Attachment[] = []

  if (Array.isArray(req?.attachment)) {
    for (let i = 0; i < req?.attachment.length; i++) {
      const attachment = req?.attachment[i]
      if (attachment) {
        const data: Attachment = {
          '@type': 'AttachmentType',
          id: attachment.id,
          attachmentType: attachment?.attachmentType || undefined,
          description: attachment?.description || undefined,
          mimeType: attachment?.mimeType || undefined,
          name: attachment?.name || undefined,
          url: attachment?.url || undefined,
          validFor: attachment?.validFor || undefined,
          redirectUrl: attachment?.redirectUrl || undefined,
          displayInfo: attachment?.displayInfo || undefined,
        }
        attachments.push(data)
      }
    }
  }

  if (req.validFor && Object.keys(req.validFor).length !== 0) {
    req.validFor.startDateTime = dayjs(req.validFor?.startDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
    req.validFor.endDateTime = dayjs(req.validFor?.endDateTime).format('YYYY-MM-DDTHH:mm:ss.SSS[Z]')
  }

  const doc: ProductLanguage = {
    id: req.id || nano.randomNanoId(),
    '@type': 'productLanguage',
    languageCode: req.languageCode,
    attachment: attachments || [],
    name: req?.name || undefined,
    version: req?.version || undefined,
    lastUpdate: req.lastUpdate || undefined,
    validFor: req?.validFor || undefined,
  }
  const result = await insertOneProductLanguage(doc)
  const productLanguage = await findProductLanguageId(result.insertedId)
  return productLanguage
}

export { createProductLanguage }
