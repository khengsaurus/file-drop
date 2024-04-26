package controllers

type ResourceInfo struct {
	FileName   string `json:"fileName"`
	Key        string `json:"key"`
	UploadedAt int64  `json:"uploadedAt"`
	Url        string `json:"url"`
}
