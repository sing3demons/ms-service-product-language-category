import { ValidFor } from './index.js'

export interface ProductLanguage {
  '@type'?: string
  id: string
  href ?: string
  languageCode: string
  attachment?: Attachment[]
  name?: string
  version?: string
  lastUpdate?: string
  validFor?: ValidFor
}

export interface Attachment {
  '@type': string
  id: string
  href ?: string
  attachmentType?: string
  description?: string
  mimeType?: string
  name?: string
  url?: string
  validFor?: ValidFor
  redirectUrl?: string
  displayInfo?: {
    valueType: string
    value: string[]
  }
}
