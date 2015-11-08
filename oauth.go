package main

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	OAuthConf = &oauth2.Config{
		ClientID:     "d103e51684009fc22250",
		ClientSecret: "63079590549474872a6b656aed3c4aaecb8a3efc",
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	OAuthStateString = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
)

func isLoggedIn() bool {
	_, err := os.Stat("token.json")
	return !os.IsNotExist(err)
}

func loadToken() (*oauth2.Token, error) {
	tokenFile, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	defer tokenFile.Close()

	err = json.NewDecoder(tokenFile).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, err
}

func saveToken(token *oauth2.Token) error {
	tokenFile, err := os.Create("token.json")
	if err != nil {
		return err
	}

	defer tokenFile.Close()
	return json.NewEncoder(tokenFile).Encode(token)
}
