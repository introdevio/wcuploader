package wp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type wordpressHttpClient struct {
	client *http.Client
	url    string
}

func (w *wordpressHttpClient) CreateRequest(method, endpoint string, buffer io.Reader) (*http.Request, error) {
	fullUrl, err := url.JoinPath(w.url, endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, fullUrl, buffer)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (w *wordpressHttpClient) PrepareMultipartFile(filePath string, fields map[string]string) (*bytes.Buffer, string, error) {
	buff := new(bytes.Buffer)

	mp := multipart.NewWriter(buff)

	for key, value := range fields {
		err := mp.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, "", err
	}

	file, err := os.Open(filePath)

	if err != nil {
		return nil, "", err
	}

	part, err := mp.CreateFormFile("file", "@"+filePath)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)

	err = mp.Close()

	if err != nil {
		return nil, "", err
	}

	return buff, mp.FormDataContentType(), nil

}

func (w *wordpressHttpClient) Execute(r *http.Request) ([]byte, int, error) {
	resp, err := w.client.Do(r)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var (
		status int
		body   string
	)

	status = resp.StatusCode
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	fmt.Println(status, body)
	return response, status, nil
}
