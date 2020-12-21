package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientId = "myclient"
	clientSecret = "bf6ba301-165f-4e42-b668-ad203d47ea3f"
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

	http.HandleFunc("/auth/callback", func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(response, "State is invalid.", http.StatusBadRequest)
			return 
		}

		//exchange the token sent by Keycloak for a access token
		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			http.Error(response, "Fail to exchange token.", http.StatusInternalServerError)
			return
		}
		
		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(response, "Fail on generating the ID token.", http.StatusInternalServerError)
			return
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
		if err != nil {
			http.Error(response, "Fail to get user information.", http.StatusInternalServerError)
			return
		}

		//transforma o toke para JSON
		resp := struct {
			AccessToken *oauth2.Token
			IDToken string
			UserInfo *oidc.UserInfo
		}{
			token, 
			idToken,
			userInfo,
		}
		
		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError) 
			return
		}

		response.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}