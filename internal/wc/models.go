package wc

import "github.com/introdevio/wcuploader/internal/product"

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Image struct {
	Id   int    `json:"id"`
	Src  string `json:"src,omitempty"`
	Name string `json:"name,omitempty"`
	Alt  string `json:"alt,omitempty"`
}

type Attribute struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name"`
	Position    int      `json:"position"`
	Visible     bool     `json:"visible"`
	IsVariation bool     `json:"variation"`
	Options     []string `json:"options"`
}

type Product struct {
	Id                int                  `json:"id"`
	Name              string               `json:"name,omitempty"`
	Type              string               `json:"type,omitempty"`
	Status            string               `json:"status,omitempty"`
	Featured          bool                 `json:"featured,omitempty"`
	CatalogVisibility string               `json:"catalog_visibility,omitempty"`
	Description       string               `json:"description,omitempty"`
	ShortDescription  string               `json:"short_description,omitempty"`
	Sku               string               `json:"sku"`
	RegularPrice      string               `json:"regular_price"`
	SalePrice         string               `json:"sale_price,omitempty"`
	ManageStock       bool                 `json:"manage_stock"`
	StockQuantity     int                  `json:"stockQuantity"`
	StockStatus       string               `json:"stock_status,omitempty"`
	Backorders        string               `json:"backorders,omitempty"`
	Weight            string               `json:"weight"`
	Tags              []Tag                `json:"tags,omitempty"`
	Images            []Image              `json:"images"`
	Attributes        []Attribute          `json:"attributes"`
	Variations        []*ProductVariation  `json:"-"`
	DefaultAttributes []VariationAttribute `json:"default_attributes"`
}

type VariationAttribute struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

type ProductVariation struct {
	Id            int                  `json:"id"`
	Description   string               `json:"description,omitempty"`
	Sku           string               `json:"sku"`
	RegularPrice  string               `json:"regular_price"`
	SalePrice     string               `json:"sale_price,omitempty"`
	Status        string               `json:"status,omitempty"`
	StockQuantity int                  `json:"stock_quantity"`
	Image         Image                `json:"image"`
	Attributes    []VariationAttribute `json:"attributes"`
}

func NewProductFromProduct(p product.Product) Product {
	var colors []string
	for c := range p.Colors {
		colors = append(colors, c)
	}
	var images []Image
	for _, img := range p.Images {

		images = append(images, Image{Id: img.RemoteImageId})
	}
	attrs := NewAttributeWithOptions("Color", colors)
	return Product{
		Name:             p.Sku,
		Type:             "variable",
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Sku:              p.Sku,
		RegularPrice:     p.RegularPrice,
		SalePrice:        p.SalePrice,
		StockQuantity:    1,
		StockStatus:      "instock",
		Images:           images,
		Variations:       VariationsFromProduct(&p),
		Attributes:       []Attribute{attrs},
	}
}

func VariationsFromProduct(p *product.Product) []*ProductVariation {
	var variations []*ProductVariation
	for _, v := range p.Variations {
		variation := ProductVariation{
			Sku:           p.Sku + "-" + v.Sku,
			RegularPrice:  v.RegularPrice,
			SalePrice:     v.SalePrice,
			StockQuantity: 1,
			Image: Image{
				Id: v.Image.RemoteImageId,
			},
			Attributes: []VariationAttribute{{Option: v.Name}},
		}
		variations = append(variations, &variation)
	}
	return variations
}

func NewAttributeWithOptions(name string, options []string) Attribute {
	return Attribute{
		Name:        name,
		Visible:     true,
		IsVariation: true,
		Options:     options,
	}
}

func NewVariationAttribute(id int, option string) VariationAttribute {
	return VariationAttribute{
		Id:     id,
		Option: option,
	}
}
