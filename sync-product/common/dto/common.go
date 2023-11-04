package dto

type ValidFor struct {
	StartDateTime string `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}
type Query struct {
	Name            string `json:"name" bson:"name,omitempty"`
	Limit           uint64 `json:"limit" bson:"limit,omitempty"`
	ID              string `json:"id" bson:"id,omitempty"`
	Depth           uint   `json:"depth" bson:"depth,omitempty"`
	Expand          string `json:"expand" bson:"expand,omitempty"`
	LifecycleStatus string `json:"lifecycleStatus" bson:"lifecycleStatus,omitempty"`
	Category        string `bson:"category,omitempty"`
	ProductLanguage string `bson:"productLanguage,omitempty"`
}
