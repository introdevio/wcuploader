package main

import (
	"errors"
	"flag"
	"github.com/introdevio/wcuploader/internal/chatgpt"
	"github.com/introdevio/wcuploader/internal/uploader"
	"github.com/introdevio/wcuploader/internal/wp"
	"log"
	url2 "net/url"
)

func main() {
	//TODO: fix instock
	//TODO: support multiple categories
	//TODO: support multiple attributes
	//TODO: externalize flags into toml config

	var (
		rootProductFolder string
		chatgptSecret     string
		wordpressUser     string
		wordpressKey      string
		woocommerceKey    string
		woocommerceSecret string
		baseUrl           string
	)

	flag.StringVar(&rootProductFolder, "products", "", "Root folder where products are located")
	flag.StringVar(&chatgptSecret, "gptsecret", "", "Chat GPT secret key")
	flag.StringVar(&wordpressUser, "user", "", "Wordpress user for Wordpress API")
	flag.StringVar(&wordpressKey, "wp-key", "", "Wordpress App Password for API")
	flag.StringVar(&woocommerceKey, "woocommerce-key", "", "API Key for Woocommerce API")
	flag.StringVar(&woocommerceSecret, "woocommerce-secret", "", "API Secret key for woocommerce api")
	flag.StringVar(&baseUrl, "shop-url", "", "Base url for woocommerce shop")
	flag.Parse()

	if wordpressUser == "" || wordpressKey == "" || woocommerceKey == "" || woocommerceSecret == "" || baseUrl == "" {
		log.Fatal(errors.New("need to define all flag requirements"))
	}

	loader := uploader.NewPathProductLoader(rootProductFolder)
	result, err := loader.Load()

	if err != nil {
		log.Fatal(err)
	}

	url, e := url2.JoinPath(baseUrl, "/wp-json/")
	if e != nil {
		log.Fatal(e)
	}
	wp2 := wp.NewWordpressAPI(wordpressUser, wordpressKey, woocommerceKey, woocommerceSecret, url)
	gpt := chatgpt.NewChatGptClient(chatgptSecret)
	categoryMap, e := wp2.GetCategories()

	if e != nil {
		log.Fatal(e)
	}

	for _, p := range result {
		for _, img := range p.Images {
			e := wp2.PostMedia(img)
			if e != nil {
				log.Fatal(e)
			}
		}

		e = gpt.CreateDescription(&p)
		if e != nil {
			log.Fatal(err)
		}

		e = gpt.CreateShortDescription(&p)
		if e != nil {
			log.Fatal(err)
		}
		product := wp.NewProductFromProduct(p, categoryMap)
		e = wp2.CreateProduct(product)
		if e != nil {
			log.Fatal(e)
		}
	}

	log.Println(result)
}
