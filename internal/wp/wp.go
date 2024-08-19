package wp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/introdevio/wcuploader/internal/product"
	"github.com/introdevio/wcuploader/internal/wc"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type WordpressAPI struct {
	user   string
	secret string
	url    string
	client *http.Client
}

func NewWordpressAPI(user, secret, url string) *WordpressAPI {
	return &WordpressAPI{
		user:   user,
		secret: secret,
		url:    url,
		client: &http.Client{},
	}
}

func (w *WordpressAPI) PostMedia(image *product.LocalImage) error {
	body, code, err := w.postMultipart("/wp/v2/media", image.Path)
	if err != nil {
		return err
	}

	var img wc.Image

	err = json.Unmarshal(body, &img)

	if err != nil {
		return err
	}

	image.RemoteImageId = img.Id

	if code != http.StatusCreated {
		log.Println("Not okay", code)
		return errors.New(fmt.Sprintf("Failed request. Reason %s", string(body)))
	}
	log.Printf("Successfully uploaded image: %s\n", filepath.Base(image.Path))
	return nil
}

func (w *WordpressAPI) ListMedia() {
	body, code, err := w.get("/wp/v2/media")

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body), code)
}

func (w *WordpressAPI) get(endpoint string) ([]byte, int, error) {

	u, err := url.JoinPath(w.url, endpoint)
	if err != nil {
		return nil, 0, err
	}

	r, err := w.client.Get(u)

	if err != nil {
		return nil, 0, err
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, 0, err
	}

	return body, r.StatusCode, nil
}

func (w *WordpressAPI) post(endpoint string, payload []byte, filePath string) ([]byte, int, error) {

	u, err := url.JoinPath(w.url, endpoint)

	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(payload))

	if err != nil {
		return nil, 0, err
	}

	req.SetBasicAuth(w.user, w.secret)

	if filePath != "" {
		req.Header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filePath))
		req.Header.Add("Content-Type", "image/jpg")
	}

	r, err := w.client.Do(req)

	if err != nil {
		return nil, 0, err
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, 0, err
	}

	return body, r.StatusCode, nil

}

func (w *WordpressAPI) postMultipart(endpoint string, media string) ([]byte, int, error) {
	reqBody := new(bytes.Buffer)
	mp := multipart.NewWriter(reqBody)
	err := mp.WriteField("title", filepath.Base(media))

	if err != nil {
		return nil, 0, err
	}

	file, err := os.Open(media)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()
	part, err := mp.CreateFormFile("file", "@"+media)
	if err != nil {

	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, 0, err
	}

	mp.Close()

	u, err := url.JoinPath(w.url, endpoint)
	if err != nil {

	}
	req, err := http.NewRequest("POST", u, reqBody)

	if err != nil {
		return nil, 0, err
	}

	req.SetBasicAuth(w.user, w.secret)
	req.Header.Add("Content-Type", mp.FormDataContentType())
	resp, err := w.client.Do(req)

	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return b, resp.StatusCode, nil

}
