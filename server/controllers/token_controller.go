package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/khengsaurus/file-drop/server/consts"
	"github.com/khengsaurus/file-drop/server/utils"
)

type tokenInfo struct {
	Token string `json:"token"`
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-> GetToken")

	token := utils.GenerateHash()

	cookie := http.Cookie{
		Name:     consts.ClientCookieName,
		Value:    token,
		Domain:   os.Getenv("CLIENT_DOMAIN"),
		Path:     "/",
		MaxAge:   2147483647,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, &cookie)
	Json200(&tokenInfo{Token: token}, w)
}
