package wc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/introdevio/wcuploader/internal/product"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WoocommerceAPI struct {
	apiKey    string
	secretKey string
	url       string
	client    *http.Client
}

func NewWoocommerceApi(apiKey, secretKey, url string) *WoocommerceAPI {
	return &WoocommerceAPI{
		apiKey:    apiKey,
		secretKey: secretKey,
		url:       url,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (wc *WoocommerceAPI) GetProductById(id string) {
	body, _, err := wc.get("products/" + id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(body)
}

func (wc *WoocommerceAPI) CreateProduct(p product.Product) {

	pr, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error", err)
	}

	rs, _, err := wc.post("products", pr)

	if err != nil {
		fmt.Println(err)
		return
	}

	var response product.Product
	err = json.Unmarshal(rs, &response)

	if err != nil {
		fmt.Println(err)
		return
	}

	createMap := make(map[string][]product.Variation)

	createMap["create"] = p.Children

	pr, err = json.Marshal(createMap)
	if err != nil {
		fmt.Println("Error", err)
	}

	r, statusCode, _ := wc.post(fmt.Sprintf("products/%d/variations", response.Id), pr)
	fmt.Println(string(r), statusCode)
}

func (wc *WoocommerceAPI) get(endpoint string) ([]byte, int, error) {
	path, err := url.JoinPath(wc.url, endpoint)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add("Content-Type", `application/json`)
	req.SetBasicAuth(wc.apiKey, wc.secretKey)

	resp, err := wc.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, resp.StatusCode, nil
}

func (wc *WoocommerceAPI) post(endpoint string, reqBody []byte) ([]byte, int, error) {
	path, err := url.JoinPath(wc.url, endpoint)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Content-Type", `application/json`)
	req.SetBasicAuth(wc.apiKey, wc.secretKey)

	resp, err := wc.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, 0, err
	}

	var (
		status int
		body   string
	)

	status = resp.StatusCode
	r, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, 0, err
	}

	fmt.Println(status, body)
	return r, status, nil
}
