import { ValidFor } from './index.js'
import { Product } from './product.js'

export interface Category {
  '@type': string
  id: string
  name: string
  version?: string
  lastUpdate?: string
  validFor?: ValidFor
  products?: Product[]
}
