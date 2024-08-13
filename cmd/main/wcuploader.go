package main

import (
	"github.com/introdevio/wcuploader/internal/uploader"
	wc2 "github.com/introdevio/wcuploader/internal/wc"
)

func main() {

	// upload to s3 bucket

	//s3bucket := "arn:aws:s3:::scottisheyewear"
	rootFolder := "/Users/daniel/Scottish/fotos-web/product/test"

	loader := uploader.NewPathProductLoader(rootFolder)
	result, _ := loader.Load()

	wc := wc2.NewWoocommerceApi(wcKey, wcSecret, url)
	//result[0].LoadProductColorVariations()
	wc.CreateProduct(result[0])
	wc.GetProductById("1629")
}
