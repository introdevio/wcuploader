package main

import (
	"github.com/introdevio/wcuploader/internal/uploader"
	wc2 "github.com/introdevio/wcuploader/internal/wc"
)

func main() {

	// upload to s3 bucket
	//TODO: remove duplicates from children
	//TODO: fix color append to whole set of attributes
	//TODO: make internal objects and API objects
	//TODO: map between internal and api objects
	//TODO: load images while loading the colors and such

	loader := uploader.NewPathProductLoader(rootFolder)
	result, _ := loader.Load()

	wc := wc2.NewWoocommerceApi(wcKey, wcSecret, url)
	//result[0].LoadProductColorVariations()
	wc.CreateProduct(result[0])
	wc.GetProductById("1629")
}
