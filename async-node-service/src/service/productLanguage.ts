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
          href: `/attachment/${attachment.id}`,
          attachmentType: attachment?.attachmentType,
          description: attachment?.description,
          mimeType: attachment?.mimeType,
          name: attachment?.name,
          url: attachment?.url,
          validFor: attachment?.validFor,
          redirectUrl: attachment?.redirectUrl,
          displayInfo: attachment?.displayInfo,
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
    name: req?.name,
    version: req?.version,
    lastUpdate: req.lastUpdate,
    validFor: req?.validFor,
  }
  const result = await insertOneProductLanguage(doc)
  const productLanguage = await findProductLanguageId(result.insertedId)
  return productLanguage
}

export { createProductLanguage }
