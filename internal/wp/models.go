package wp

type Media struct {
	Id            int    `json:"id,omitempty"`
	Link          string `json:"link,omitempty"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	CommentStatus string `json:"comment_status"`
	PingStatus    string `json:"ping_status"`
	AltText       string `json:"alt_text,omitempty"`
	Description   string `json:"description"`
	SourceUrl     string `json:"source_url,omitempty"`
}

type MediaResponse struct {
	Id   int    `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
}
