package main

import (
	"fmt"
	"github.com/introdevio/wcuploader/internal/uploader"
	"log"
)

func main() {

	// upload to s3 bucket

	//s3bucket := "arn:aws:s3:::scottisheyewear"
	rootFolder := "/Users/daniel/Scottish/fotos-web/product"
	// load csv
	csvPath := rootFolder + "/products.csv"

	loader := uploader.NewProductLoader(csvPath, rootFolder)

	result, err := loader.LoadFromCsv(true)

	fmt.Printf("%+v", result)

	if err != nil {
		log.Fatalln(err)
	}

}
