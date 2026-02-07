package models

type SamokatResponseDTO struct {
	Categories []CategoryDTO `json:"categories"`
}

type CategoryDTO struct {
	Name     string       `json:"name"`
	Products []ProductDTO `json:"products,omitempty"`
}

type ProductDTO struct {
	Name   string     `json:"name"`
	Slug   string     `json:"slug,omitempty"`
	Media  []MediaDTO `json:"media,omitempty"`
	Prices *PricesDTO `json:"prices,omitempty"`
}

type MediaDTO struct {
	URL string `json:"url"`
}

type PricesDTO struct {
	Current int `json:"current"`
}
