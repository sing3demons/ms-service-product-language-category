import { Category } from './category.js'
import { ValidFor } from './index.js'
import { ProductLanguage } from './productLanguage.js'

export interface Product {
  '@type'?: string
  id: string
  href?: string
  lastUpdate?: string
  lifecycleStatus?: string
  name?: string
  version?: string
  category?: Category[]
  price?: ProductPrice[]
  description?: string
  supportingLanguage?: ProductLanguage[]
  validFor?: ValidFor
}

export interface ProductPrice {
  id?: string
  name?: string
}
