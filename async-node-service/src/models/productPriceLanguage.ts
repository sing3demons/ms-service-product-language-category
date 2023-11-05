export interface IProductPriceLanguageDTO {
  '@type'?: string
  id: string
  languageCode: string
  lastUpdate?: Date
  name?: string
  version?: string
  price?: IPrice
  validFor?: ITimePeriod
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
