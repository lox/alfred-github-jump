package jump

import (
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
	ServerAddress    = "http://127.0.0.1:7024"
	ServerBind       = ":7024"
)
