package gcp

import (
	"os"

	"golang.org/x/oauth2"
)

func getStaticTokenSource() (oauth2.TokenSource, bool) {
	v, ok := os.LookupEnv("GOOGLE_OAUTH_ACCESS_TOKEN")
	if ok {
		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: v}), true
	}

	return nil, false
}
