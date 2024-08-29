package wp

import (
	"github.com/introdevio/wcuploader/internal"
)

type Media struct {
	Id            int    `json:"id,omitempty"`
	Link          string `json:"link,omitempty"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	CommentStatus string `json:"comment_status"`
	PingStatus    string `json:"ping_status"`
	AltText       string `json:"alt_text,omitempty"`
	Description   string `json:"description"`
	SourceUrl     string `json:"source_url,omitempty"`
}

type MediaResponse struct {
	Id   int    `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
}

var CategoryMap = map[string]int{"damas": 22, "caballeros": 23, "sobrelentes": 22}

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
	Categories        []Category           `json:"categories"`
	Tags              []Tag                `json:"tags,omitempty"`
	Images            []Image              `json:"images"`
	Attributes        []Attribute          `json:"attributes"`
	Variations        []*ProductVariation  `json:"-"`
	DefaultAttributes []VariationAttribute `json:"default_attributes"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type VariationAttribute struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

type ProductVariation struct {
	Id            int                  `json:"id,omitempty"`
	Description   string               `json:"description,omitempty"`
	Sku           string               `json:"sku"`
	RegularPrice  string               `json:"regular_price"`
	SalePrice     string               `json:"sale_price,omitempty"`
	Status        string               `json:"status,omitempty"`
	StockQuantity int                  `json:"stock_quantity"`
	Image         Image                `json:"image"`
	Attributes    []VariationAttribute `json:"attributes"`
}

func NewProductFromProduct(p internal.Product) Product {
	var colors []string
	for c := range p.Colors {
		colors = append(colors, c)
	}
	var images []Image
	for _, img := range p.Images {

		images = append(images, Image{Id: img.RemoteImageId})
	}
	attrs := NewAttributeWithOptions("Color", colors)
	var categories []Category
	for _, category := range p.Categories {
		categories = append(categories, Category{
			Id: CategoryMap[category],
		})
	}
	return Product{
		Name:             p.Sku,
		Type:             "variable",
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Sku:              p.Sku,
		RegularPrice:     p.RegularPrice,
		SalePrice:        p.SalePrice,
		StockQuantity:    1,
		ManageStock:      true,
		StockStatus:      "instock",
		Images:           images,
		Categories:       categories,
		Variations:       VariationsFromProduct(&p),
		Attributes:       []Attribute{attrs},
	}
}

func VariationsFromProduct(p *internal.Product) []*ProductVariation {
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
			Attributes: []VariationAttribute{{Option: v.Name, Name: "Color"}},
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
