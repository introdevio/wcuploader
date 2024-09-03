package wp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/introdevio/wcuploader/internal"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type WordpressAPI struct {
	user              string
	password          string
	woocommerceKey    string
	woocommerceSecret string
	client            *wordpressHttpClient
}

func NewWordpressAPI(user, password, woocommerceKey, woocommerceSecret, apiUrl string) *WordpressAPI {
	return &WordpressAPI{
		user:              user,
		password:          password,
		woocommerceKey:    woocommerceKey,
		woocommerceSecret: woocommerceSecret,
		client: &wordpressHttpClient{
			client: &http.Client{},
			url:    apiUrl,
		},
	}
}

func (w *WordpressAPI) GetCategories() (map[string]int, error) {
	req, err := w.client.CreateRequest(http.MethodGet, "wc/v3/products/categories", nil)

	if err != nil {
		return nil, err
	}
	req.URL.Query().Add("per_page", "100")
	req.SetBasicAuth(w.woocommerceKey, w.woocommerceSecret)

	resp, _, err := w.client.Execute(req)

	if err != nil {
		return nil, err
	}

	var categories []Category
	err = json.Unmarshal(resp, &categories)

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, category := range categories {
		result[strings.ToLower(category.Name)] = category.Id
	}

	return result, nil
}

func (w *WordpressAPI) GetProductById(id string) (*Product, error) {
	req, err := w.client.CreateRequest(http.MethodGet, "wc/v3/products/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(w.woocommerceKey, w.woocommerceSecret)
	response, status, err := w.client.Execute(req)
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, errors.New(fmt.Sprintf("product with id: %s not found", id))
	}
	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Something did not work as expected when fetching product %s", string(response)))
	}

	var p Product
	err = json.Unmarshal(response, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (w *WordpressAPI) CreateProduct(p Product) error {

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := w.client.CreateRequest(http.MethodPost, "wc/v3/products", bytes.NewBuffer(payload))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", `application/json`)
	req.SetBasicAuth(w.woocommerceKey, w.woocommerceSecret)

	response, status, err := w.client.Execute(req)

	if err != nil {
		return err
	}

	if status > http.StatusAccepted {
		return errors.New(fmt.Sprintf("Something did not work as expected when creating product %s", string(response)))
	}

	var created Product
	err = json.Unmarshal(response, &created)

	if err != nil {
		return err
	}

	log.Printf("Created Product with ID: %d\n", created.Id)

	for _, v := range p.Variations {
		for _, attr := range v.Attributes {
			attr.Id = created.Attributes[0].Id
		}
	}

	variations := make(map[string][]*ProductVariation)
	variations["create"] = p.Variations
	payload, err = json.Marshal(variations)

	req, err = w.client.CreateRequest(http.MethodPost, fmt.Sprintf("wc/v3/products/%d/variations/batch", created.Id), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", `application/json`)
	req.SetBasicAuth(w.woocommerceKey, w.woocommerceSecret)
	response, status, err = w.client.Execute(req)
	if err != nil {
		return err
	}

	if status > http.StatusAccepted {
		return errors.New(fmt.Sprintf("Something did not work as expected when creating product %s", string(response)))
	}

	var result map[string][]ProductVariation
	err = json.Unmarshal(response, &result)

	if err != nil {
		return err
	}
	return nil
}

func (w *WordpressAPI) PostMedia(image *internal.LocalImage) error {
	fields := make(map[string]string)
	fields["title"] = filepath.Base(image.Path)
	mp, contentType, err := w.client.PrepareMultipartFile(image.Path, fields)
	if err != nil {
		return err
	}

	req, err := w.client.CreateRequest(http.MethodPost, "/wp/v2/media", mp)

	if err != nil {
		return err
	}
	req.SetBasicAuth(w.user, w.password)
	req.Header.Add("Content-Type", contentType)
	response, status, err := w.client.Execute(req)

	if status != http.StatusCreated {
		return errors.New(fmt.Sprintf("Failed to upload media to wp. Reason %s", string(response)))
	}

	var img MediaResponse
	err = json.Unmarshal(response, &img)

	if err != nil {
		return err
	}

	image.RemoteImageId = img.Id
	image.RemoteUrl = img.Link

	log.Printf("Successfully uploaded image: %s\n", filepath.Base(image.Path))
	return nil
}

func (w *WordpressAPI) ListMedia() ([]MediaResponse, error) {
	req, err := w.client.CreateRequest(http.MethodGet, "/wp/v2/media", nil)

	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(w.user, w.password)
	response, status, err := w.client.Execute(req)

	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Failed to get media from wp. Reason %s", string(response)))
	}
	var media []MediaResponse
	err = json.Unmarshal(response, &media)
	return media, nil
}
