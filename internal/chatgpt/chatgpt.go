package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/introdevio/wcuploader/internal"
	"io"
	"net/http"
	"net/url"
)

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type MessageContent struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageUrl ImageUrl `json:"image_url,omitempty"`
}

type ImageUrl struct {
	Url string `json:"url"`
}

type ChatResponse struct {
	Id                string           `json:"id"`
	Object            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	SystemFingerprint string           `json:"system_fingerprint"`
	Choices           []ResponseChoice `json:"choices"`
	Usage             Usage            `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ResponseChoice struct {
	Index   int             `json:"index"`
	Message ResponseMessage `json:"message"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	apiSecret string
	client    *http.Client
	url       string
}

func NewChatRequest(messages []Message) ChatRequest {
	return ChatRequest{
		Model:    "gpt-4o",
		Messages: messages,
	}
}

func NewMessageWithImage(prompt, imageUrl string) Message {
	p := MessageContent{
		Type: "text",
		Text: prompt,
	}

	img := MessageContent{
		Type:     "image_url",
		ImageUrl: ImageUrl{imageUrl},
	}
	return Message{
		Role:    "user",
		Content: []MessageContent{p, img},
	}
}

func NewChatGptClient(apiSecret string) *Client {
	return &Client{
		apiSecret: apiSecret,
		client:    &http.Client{},
		url:       "https://api.openai.com/v1/",
	}
}

func (c *Client) CreateShortDescription(p *internal.Product) error {
	prompt := fmt.Sprintf("Crea una description corta para este lente de categoria %s modelo %s sin incluir dimensiones y sin formato", p.Categories[0], p.Sku)
	msg := NewMessageWithImage(prompt, p.Images[0].RemoteUrl)
	chat := NewChatRequest([]Message{msg})
	u, err := url.JoinPath(c.url, "chat/completions")
	if err != nil {
		return err
	}
	body, err := json.Marshal(chat)
	if err != nil {
		return err
	}
	req, err := c.createRequest("POST", u, body)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result ChatResponse

	err = json.Unmarshal(data, &result)

	if err != nil {
		return err
	}

	desc := result.Choices[0].Message.Content
	p.ShortDescription = desc
	fmt.Printf("Created Short Description. Used %d tokens\n", result.Usage.TotalTokens)
	return nil
}

func (c *Client) CreateDescription(p *internal.Product) error {
	prompt := fmt.Sprintf("Crea una description detallada para este lente de categoria %s modelo %s sin incluir dimensiones y formatealo en simple html for a wp page and do not add markup code wrapper", p.Categories[0], p.Sku)
	msg := NewMessageWithImage(prompt, p.Images[0].RemoteUrl)
	chat := NewChatRequest([]Message{msg})
	u, err := url.JoinPath(c.url, "chat/completions")
	if err != nil {
		return err
	}
	body, err := json.Marshal(chat)
	if err != nil {
		return err
	}
	req, err := c.createRequest("POST", u, body)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result ChatResponse

	err = json.Unmarshal(data, &result)

	if err != nil {
		return err
	}

	desc := result.Choices[0].Message.Content
	p.Description = desc
	fmt.Printf("Created Description. Used %d tokens\n", result.Usage.TotalTokens)
	return nil
}

func (c *Client) createRequest(method, endpoint string, body []byte) (*http.Request, error) {
	r, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiSecret))
	r.Header.Add("Content-Type", "application/json")

	return r, nil
}
