export interface IProductPriceDTO {
  '@type'?: string
  id?: string
  lastUpdate?: Date
  lifecycleStatus?: string
  name?: string
  version?: string
  price?: IPrice
  validFor?: ITimePeriod
  popRelationship?: IPopRelationship[]
  TORO_supportingLanguage?: TORO_SupportingLanguage[]
  deleteDate?: Date
}

export interface IPrice {
  unit?: string
  value?: number
}

export interface ITimePeriod {
  endDateTime?: Date
  startDateTime?: Date
}

export interface IPopRelationship {
  id?: string
  name?: string
}

export interface TORO_SupportingLanguage {
  id?: string
  languageCode?: string
  referredType?: string
}
