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
