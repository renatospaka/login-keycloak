package main

import (
	"context"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientId = "myclient"
	clientSecret = "fc81e341-2415-4a06-9acd-52a9d5a3b4f8"
)

func main() {
	ctx := context.Background()	

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/myrealm")
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID: clientId,
    ClientSecret: clientSecret,
    Endpoint: provider.Endpoint(),
    RedirectURL: "http://localhost:8081/auth/callback",
    Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123UmaSenhaQualquer!"

http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, config.AuthCodeURL(state), http.StatusFound)
})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
