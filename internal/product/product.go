package product

import (
	"fmt"
	"path/filepath"
)

type Product struct {
	Id               int          `json:"id"`
	PostTitle        string       `csv:"post_title" json:"name"`
	ParentId         int          `json:"parent_id"`
	Description      string       `csv:"description" json:"description"`
	ShortDescription string       `csv:"short_description" json:"short_description"`
	ProductType      string       `csv:"product_type" json:"type"`
	Tags             []string     `csv:"tags" json:"-"`
	Categories       []Category   `csv:"categories" json:"categories"`
	Sku              string       `csv:"sku" json:"sku"`
	Stock            bool         `csv:"stock" json:"stock"`
	StockStatus      string       `csv:"stock_status" json:"stock_status"`
	Status           string       `csv:"status" json:"status"`
	RegularPrice     string       `csv:"regular_price" json:"regular_price"`
	SalePrice        string       `csv:"sale_price" json:"sale_price"`
	Images           []Image      `csv:"images" json:"images"`
	Children         []Variation  `csv:"-" json:"-"`
	Attributes       []*Attribute `json:"attributes"`
}

type Variation struct {
	Id           int                `json:"id"`
	Description  string             `json:"description"`
	Sku          string             `json:"sku"`
	RegularPrice string             `json:"regular_price"`
	SalePrice    string             `json:"sale_price"`
	Status       string             `json:"status"`
	ManageStock  bool               `json:"manage_stock"`
	StockStatus  string             `json:"stock_status"`
	Image        Image              `json:"image"`
	Attribute    VariationAttribute `json:"attribute"`
}

type Image struct {
	Src  string `json:"src"`
	Name string `json:"name"`
	Alt  string `json:"alt"`
}

type Attribute struct {
	Name      string   `json:"name"`
	Visible   bool     `json:"visible"`
	Variation bool     `json:"variation"`
	Options   []string `json:"options"`
}

type VariationAttribute struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Visible   bool   `json:"visible"`
	Variation bool   `json:"variation"`
	Option    string `json:"option"`
}

type Category struct {
	Id int `json:"id"`
}

func NewProductImage(src, name string) Image {
	return Image{
		Src:  src,
		Name: name,
	}
}

func (p *Product) LoadMedia(root string) {
	dir := filepath.Join(root, p.Sku, "*.jpg")
	media, err := filepath.Glob(dir)
	if err != nil {
		fmt.Println("Error reading files")
		return
	}
	for _, media := range media {
		img := Image{
			Src:  media,
			Name: p.Sku,
		}
		p.Images = append(p.Images, img)
	}
}
