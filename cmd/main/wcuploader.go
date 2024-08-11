package main

import (
	"fmt"
	"github.com/introdevio/wcuploader/internal/uploader"
)

func main() {

	// upload to s3 bucket

	//s3bucket := "arn:aws:s3:::scottisheyewear"
	rootFolder := "/Users/daniel/Scottish/fotos-web/product/test"

	loader := uploader.NewPathProductLoader(rootFolder)
	result, _ := loader.Load()

	fmt.Printf("%+v", result)

}
