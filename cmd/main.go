package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var (
	conf *oauth2.Config
	ctx  context.Context
)

func doExchange(token string) (string, error) {
	d := url.Values{}
	d.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	d.Add("assertion", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.googleapis.com/oauth2/v4/token", strings.NewReader(d.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	queryParts, _ := url.ParseQuery(r.URL.RawQuery)

	// Use the authorization code that is pushed to the redirect
	// URL.
	code := queryParts["code"][0]
	log.Printf("code: %s\n", code)

	// Exchange will do the handshake to retrieve the initial access token.
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Token: %s", tok)

	tox, err := doExchange(tok.AccessToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Token: %s", tox)
	// The HTTP Client returned by conf.Client will refresh the token as necessary.
	// client := conf.Client(ctx, tok)

	// resp, err := client.Get(Google)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println(color.CyanString("Authentication successful"))
	// }
	// defer resp.Body.Close()

	payload, err := idtoken.Validate(context.Background(), tok.AccessToken, "your google client id")
	if err != nil {
		panic(err)
	}
	fmt.Print(payload.Claims)

	fmt.Fprintf(w, `
	<html>
	<body>
	<p><strong>Success!</strong></p>
	<p>You are authenticated and can now return to the CLI.</p>
	</body>
	</html>
	`)
}

func main() {
	ctx = context.Background()
	conf = &oauth2.Config{
		ClientID:     "1021096195684-2mj1j3tmiug66qfv5jv2eaj454a2i703.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid"},
		Endpoint: google.Endpoint,
		// RedirectURL: "http://localhost:3000/oauth/callback",
		RedirectURL: "http://localhost:3000/api/auth/callback/google",
	}

	// add transport for self-signed certificate to context
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Println(color.CyanString("You will now be taken to your browser for authentication"))
	time.Sleep(1 * time.Second)
	open.Run(url)
	time.Sleep(1 * time.Second)
	log.Printf("Authentication URL: %s\n", url)

	// http.HandleFunc("/oauth/callback", callbackHandler)
	http.HandleFunc("/api/auth/callback/google", callbackHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
