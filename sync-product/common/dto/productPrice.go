package dto

type ProductPrice struct {
	Type               string               `json:"@type,omitempty" bson:"@type,omitempty"`
	Href               string               `json:"href,omitempty" bson:"href,omitempty"`
	Id                 string               `json:"id,omitempty" bson:"id,omitempty"`
	LifecycleStatus    string               `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	LastUpdate         string               `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Name               string               `json:"name,omitempty" bson:"name,omitempty"`
	Version            string               `json:"version,omitempty" bson:"version,omitempty"`
	Price              Price                `json:"price,omitempty" bson:"price,omitempty"`
	ValidFor           ValidFor             `json:"validFor,omitempty" bson:"validFor,omitempty"`
	PopRelationship    []PopRelationship    `json:"popRelationship,omitempty" bson:"popRelationship,omitempty"`
	SupportingLanguage []SupportingLanguage `json:"supportingLanguage,omitempty" bson:"supportingLanguage,omitempty"`
}

type Price struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}

type PopRelationship struct {
	Id   string `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type SupportingLanguage struct {
	Id           string `json:"id" bson:"id"`
	LanguageCode string `json:"languageCode" bson:"languageCode"`
	ReferredType string `json:"@referredType" bson:"@referredType"`
}
