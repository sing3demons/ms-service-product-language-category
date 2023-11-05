package model

type ProductLanguage struct {
	Type         string       `json:"@type" validate:"required" bson:"@type"`
	ID           string       `json:"id" validate:"required" bson:"id" binding:"required"`
	Href         string       `json:"href" bson:"href"`
	LanguageCode string       `json:"languageCode" bson:"languageCode" binding:"required"`
	Name         string       `json:"name" bson:"name"`
	Version      string       `json:"version" bson:"version"`
	LastUpdate   string       `json:"lastUpdate" bson:"lastUpdate"`
	ValidFor     *ValidFor    `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Attachment   []Attachment `json:"attachment" bson:"attachment"`
}

type Attachment struct {
	Type        string      `json:"@type" validate:"required" bson:"@type"`
	ID          string      `json:"id" validate:"required" bson:"id" binding:"required"`
	Href        string      `json:"href" bson:"href"`
	Description string      `json:"description" bson:"description"`
	MimeType    string      `json:"mimeType" bson:"mimeType"`
	Name        string      `json:"name" bson:"name"`
	Url         string      `json:"url" bson:"url"`
	ValidFor    *ValidFor   `json:"validFor,omitempty" bson:"validFor,omitempty"`
	RedirectUrl string      `json:"redirectUrl" bson:"redirectUrl"`
	DisplayInfo DisplayInfo `json:"displayInfo" bson:"displayInfo"`
}

type DisplayInfo struct {
	ValueType string   `json:"valueType" bson:"valueType"`
	Value     []string `json:"value" bson:"value"`
}

type ValidFor struct {
	StartDateTime string `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}

// type Product struct {
// 	Type       string   `json:"@type" bson:"@type,omitempty"`
// 	ID         string   `json:"id" bson:"id,omitempty"`
// 	Name       string   `json:"name" bson:"name,omitempty"`
// 	Version    string   `json:"version" bson:"version,omitempty"`
// 	LastUpdate string   `json:"lastUpdate" bson:"lastUpdate,omitempty"`
// 	ValidFor   ValidFor `json:"validFor" bson:"validFor,omitempty"`
// }

type Query struct {
	Name            string `json:"name" bson:"name,omitempty"`
	Limit           uint64 `json:"limit" bson:"limit,omitempty"`
	ID              string `json:"id" bson:"id,omitempty"`
	Depth           uint   `json:"depth" bson:"depth,omitempty"`
	Expand          string `json:"expand" bson:"expand,omitempty"`
	LifecycleStatus string `json:"lifecycleStatus" bson:"lifecycleStatus,omitempty"`
}
