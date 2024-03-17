package controllers

import (
	"fmt"
	"net/http"

	"github.com/khengsaurus/file-drop/server/utils"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> GetToken")

	Json200(&TokenInfo{Token: utils.GenerateHash()}, w)
}

// --------------- types ---------------

type TokenInfo struct {
	Token string `json:"token"`
}
