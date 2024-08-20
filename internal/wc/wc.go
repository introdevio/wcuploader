package wc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (wc *WoocommerceAPI) CreateProduct(p Product) error {

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	response, _, err := wc.post("products", payload)
	if err != nil {
		return err
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
	vPayload := make(map[string][]*ProductVariation)
	vPayload["create"] = p.Variations
	variations, err := json.Marshal(vPayload)

	vResponse, _, err := wc.post(fmt.Sprintf("products/%d/variations/batch", created.Id), variations)

	if err != nil {
		return err
	}

	var result map[string][]ProductVariation
	err = json.Unmarshal(vResponse, &result)

	if err != nil {
		return err
	}
	return nil

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
	if err != nil {
		return nil, 0, err
	}
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
