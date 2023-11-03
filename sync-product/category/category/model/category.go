package model

import "github.com/sing3demons/product.product.sync/common"

type Category struct {
	Type            string        `json:"@type,omitempty" validate:"required" bson:"@type,omitempty"`
	ID              string        `json:"id,omitempty" validate:"required" bson:"id,omitempty"`
	Href            string        `json:"href,omitempty" bson:"href,omitempty"`
	Name            string        `json:"name,omitempty" bson:"name,omitempty"`
	Version         string        `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string        `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	LifecycleStatus common.Status `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	ValidFor        *ValidFor     `json:"validFor,omitempty" bson:"validFor,omitempty"`
	Products        []ProductRef     `json:"products,omitempty" bson:"product,omitempty"`
}

type UpdateCategory struct {
	ID              string        `json:"id,omitempty" validate:"required" bson:"id,omitempty"`
	Name            string        `json:"name,omitempty" bson:"name,omitempty"`
	Version         string        `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate      string        `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	LifecycleStatus common.Status `json:"lifecycleStatus,omitempty" bson:"lifecycleStatus,omitempty"`
	Products        []ProductRef     `json:"products,omitempty" bson:"product,omitempty"`
}

type ValidFor struct {
	StartDateTime string `json:"startDateTime,omitempty" bson:"startDateTime,omitempty"`
	EndDateTime   string `json:"endDateTime,omitempty" bson:"endDateTime,omitempty"`
}

type Product struct {
	Type       string    `json:"@type" bson:"@type,omitempty"`
	ID         string    `json:"id,omitempty" bson:"id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Version    string    `json:"version,omitempty" bson:"version,omitempty"`
	LastUpdate string    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	ValidFor   *ValidFor `json:"validFor,omitempty" bson:"validFor,omitempty"`
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

type Query struct {
	Name            string `json:"name" bson:"name,omitempty"`
	Limit           uint64 `json:"limit" bson:"limit,omitempty"`
	ID              string `json:"id" bson:"id,omitempty"`
	Depth           uint   `json:"depth" bson:"depth,omitempty"`
	Expand          string `json:"expand" bson:"expand,omitempty"`
	LifecycleStatus string `json:"lifecycleStatus" bson:"lifecycleStatus,omitempty"`
}
