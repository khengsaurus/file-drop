package types

type HtmlPageImg struct {
	Title string
	Src   string
}

type ResourceInfo struct {
	FileName string `json:"fileName"`
	Key      string `json:"key"`
	Url      string `json:"url"`
}

type UrlInfo struct {
	Url string `json:"url"`
	Key string `json:"key"`
}
