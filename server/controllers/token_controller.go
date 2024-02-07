package controllers

import (
	"fmt"
	"net/http"

	"github.com/khengsaurus/file-drop/server/types"
	"github.com/khengsaurus/file-drop/server/utils"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> GetToken")

	utils.Json200(&types.TokenInfo{Token: utils.GenerateHash()}, w)
}
