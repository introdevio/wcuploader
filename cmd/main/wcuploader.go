package main

import (
	"errors"
	"flag"
	"github.com/introdevio/wcuploader/internal/uploader"
	"github.com/introdevio/wcuploader/internal/wc"
	"github.com/introdevio/wcuploader/internal/wp"
	"log"
)

func main() {

	var (
		rootProductFolder string
		wordpressUser     string
		wordpressKey      string
		woocommerceKey    string
		woocommerceSecret string
		baseUrl           string
	)

	flag.StringVar(&rootProductFolder, "products", "", "Root folder where products are located")
	flag.StringVar(&wordpressUser, "user", "", "Wordpress user for Wordpress API")
	flag.StringVar(&wordpressKey, "wordpress-key", "", "Wordpress App Password for API")
	flag.StringVar(&woocommerceKey, "woocommerce-key", "", "API Key for Woocommerce API")
	flag.StringVar(&woocommerceSecret, "woocommerce-secret", "", "API Secret key for woocommerce api")
	flag.StringVar(&baseUrl, "shop-url", "", "Base url for woocommerce shop")
	flag.Parse()

	if wordpressUser == "" || wordpressKey == "" || woocommerceKey == "" || woocommerceSecret == "" {
		log.Fatal(errors.New("need to define all flag requirements"))
	}

	loader := uploader.NewPathProductLoader(rootProductFolder)
	result, err := loader.Load()

	if err != nil {
		log.Fatal(err)
	}

	url := baseUrl + "/wp-json/"
	wcUrl := url + "wc/v3"
	wp2 := wp.NewWordpressAPI(wordpressUser, wordpressKey, url)
	wc2 := wc.NewWoocommerceApi(woocommerceKey, woocommerceSecret, wcUrl)

	for _, img := range result[0].Images {
		e := wp2.PostMedia(img)
		if e != nil {
			log.Fatal(e)
		}
	}
	product := wc.NewProductFromProduct(result[0])
	e := wc2.CreateProduct(product)
	if e != nil {
		log.Fatal(e)
	}

	log.Println(result)
}
