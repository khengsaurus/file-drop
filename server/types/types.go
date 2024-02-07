package types

type HtmlPageImg struct {
	Title string
	Src   string
}

type ResourceInfo struct {
	FileName   string `json:"fileName"`
	Key        string `json:"key"`
	UploadedAt int64  `json:"uploadedAt"`
	Url        string `json:"url"`
}

type UrlInfo struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

type TokenInfo struct {
	Token string `json:"token"`
}
