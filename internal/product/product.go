package product

import (
	"fmt"
	"path/filepath"
)

type Product struct {
	PostTitle        string   `csv:"post_title"`
	Description      string   `csv:"description"`
	ShortDescription string   `csv:"short_description"`
	Tags             []string `csv:"tags"`
	Categories       []string `csv:"categories"`
	Sku              string   `csv:"sku"`
	Stock            int      `csv:"stock"`
	StockStatus      string   `csv:"stock_status"`
	Backorders       string   `csv:"backorders"`
	RegularPrice     float32  `csv:"regular_price"`
	SalePrice        float32  `csv:"sale_price"`
	Weight           float32  `csv:"weight"`
	Length           float32  `csv:"length"`
	Width            float32  `csv:"width"`
	Height           float32  `csv:"height"`
	Images           []string `csv:"images"`
}

func (p *Product) LoadMedia(root string) {
	dir := filepath.Join(root, p.Sku, "*.jpg")
	media, err := filepath.Glob(dir)
	if err != nil {
		fmt.Println("Error reading files")
		return
	}
	for _, media := range media {
		p.Images = append(p.Images, media)
	}
}
