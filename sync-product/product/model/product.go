package model

import "time"

type ProductLanguage struct {
	Type         string       `json:"@type" validate:"required" bson:"@type"`
	ID           string       `json:"id,omitempty" validate:"required" bson:"id,omitempty" binding:"required"`
	Href         string       `json:"href,omitempty" bson:"href,omitempty"`
	LanguageCode string       `json:"languageCode,omitempty" bson:"languageCode,omitempty" binding:"required"`
	ReferredType string       `json:"@referredType,omitempty" bson:"@referredType,omitempty"`
	Name         string       `json:"name,omitempty" bson:"name,omitempty"`
	Version      string       `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate   string       `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor     *ValidFor    `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Attachment   []Attachment `json:"attachment,omitempty" bson:"attachment,omitempty"`
}

type Attachment struct {
	Type        string      `json:"@type" validate:"required" bson:"@type"`
	ID          string      `json:"id" validate:"required" bson:"id" binding:"required"`
	Href        string      `json:"href,omitempty" bson:"href,omitempty"`
	Description string      `json:"description,omitempty" bson:"description,omitempty"`
	MimeType    string      `json:"mimeType,omitempty" bson:"mimeType,omitempty"`
	Name        string      `json:"name,omitempty" bson:"name,omitempty"`
	Url         string      `json:"url,omitempty" bson:"url,omitempty"`
	ValidFor    *ValidFor   `json:"validFor,omitempty" bson:"validFor,omitempty"`
	RedirectUrl string      `json:"redirectUrl,omitempty" bson:"redirectUrl,omitempty"`
	DisplayInfo DisplayInfo `json:"displayInfo,omitempty" bson:"displayInfo,omitempty"`
}

type DisplayInfo struct {
	ValueType string   `json:"valueType,omitempty" bson:"valueType,omitempty"`
	Value     []string `json:"value,omitempty" bson:"value,omitempty"`
}

type ValidForDate struct {
	StartDateTime time.Time `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   time.Time `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}

type ResponseDataWithTOtal struct {
	Total    int64      `json:"total" bson:"total"`
	Products []Products `json:"products" bson:"products"`
}

type Category struct {
	Type       string    `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID         string    `json:"id" validate:"required" bson:"id"`
	Href       string    `json:"href,omitempty" bson:"href,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Version    string    `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate string    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor   *ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
	LifecycleStatus string `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	Products   []Products `json:"products,omitempty" bson:"products,omitempty"`
}

type Products struct {
	Type               string            `json:"@type,omitempty" bson:"@type,omitempty"`
	ID                 string            `json:"id" bson:"id,omitempty"`
	Name               string            `json:"name,omitempty" bson:"name,omitempty"`
	Href               string            `json:"href,omitempty" bson:"href,omitempty"`
	LifecycleStatus    string            `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	Version            string            `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate         string            `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor           *ValidFor         `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Category           []Category        `json:"category,omitempty" bson:"category,omitempty"`
	SupportingLanguage []ProductLanguage `json:"supportingLanguage,omitempty" bson:"supportingLanguage,omitempty"`
}

type ValidFor struct {
	StartDateTime string `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}
