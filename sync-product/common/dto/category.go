package dto

type Category struct {
	Type            string     `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID              string     `json:"id,omitempty" validate:"required" bson:"id,omitempty"`
	Href            string     `json:"href,omitempty" bson:"href,omitempty"`
	Name            string     `json:"name,omitempty" bson:"name,omitempty"`
	Version         string     `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string     `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	LifecycleStatus string     `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	ValidFor        *ValidFor  `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Products        []Products `json:"products,omitempty" bson:"products,omitempty"`
}

type CategoryProducts struct {
	Type            string     `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID              string     `json:"id,omitempty" validate:"required" bson:"id,omitempty"`
	Href            string     `json:"href,omitempty" bson:"href,omitempty"`
	Name            string     `json:"name,omitempty" bson:"name,omitempty"`
	Version         string     `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string     `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	LifecycleStatus string     `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	ValidFor        *ValidFor  `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Products        []Products `json:"products,omitempty" bson:"products,omitempty"`
}

type ProductRef struct {
	Type       string    `json:"@type" bson:"@type,omitempty"`
	ID         string    `json:"id,omitempty" bson:"id,omitempty"`
	Href       string    `json:"href,omitempty" bson:"href,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Version    string    `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate string    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor   *ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
}
