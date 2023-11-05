package dto

type ProductPriceLanguage struct {
	Type         string   `json:"@type,omitempty" bson:"@type,omitempty"`
	Id           string   `json:"id,omitempty" bson:"id,omitempty"`
	Href         string   `json:"href,omitempty" bson:"href,omitempty"`
	LanguageCode string   `json:"languageCode,omitempty" bson:"languageCode,omitempty"`
	LastUpdate   string   `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Version      string   `json:"version,omitempty" bson:"version,omitempty"`
	Price        IPrice   `json:"price,omitempty" bson:"price,omitempty"`
	ValidFor     ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
}

type IPrice struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}
