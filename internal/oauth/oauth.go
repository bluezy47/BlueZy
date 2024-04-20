package oauth

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"net/url"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

type OAuth struct {
	clientID string // private
	clientSecret string // private
	redirectURL string // private
}

func NewOAuth() (*OAuth, error ) {
	var clientID string = os.Getenv("CLIENT_ID")
	var clientSecret string = os.Getenv("CLIENT_SECRET")
	var redirectURL string = os.Getenv("REDIRECT_URL")
	//
	if clientID == "" || clientSecret == "" {
		log.Println("::[OAuth][NewOAuth] Client ID or Client Secret is not set ::");
		return nil, fmt.Errorf("client id or client secret is not set");
	}
	//
	if redirectURL == "" {
		log.Println("::[OAuth][NewOAuth] Redirect URL is not set ::");
		return nil, fmt.Errorf("redirect url is not set");
	}
	if redirectURL[:4] != "http" {
		log.Println("::[OAuth][NewOAuth] Redirect URL is not valid ::");
		return nil, fmt.Errorf("redirect url is not valid");
	}
	// all ok...
	return &OAuth{
		clientID: clientID,
		clientSecret: clientSecret,
		redirectURL: redirectURL,
	}, nil
}


//	RedirectGoogleLogin redirects the user to the google oauth consent screen
func (o *OAuth) RedirectGoogleLogin(c fiber.Ctx) error {
	if o.clientID == "" || o.clientSecret == "" {
		log.Println("::[OAuth][RedirectGoogleLogin] Client ID or Client Secret is not set ::");
		// send a 500 status code. 
		return c.SendStatus(http.StatusInternalServerError)
	}
	redirectURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=https://www.googleapis.com/auth/userinfo.email&prompt=select_account", o.clientID, o.redirectURL);
	//
	return c.Redirect().To(redirectURL)
}


// HandleInitCode handles the code from the google oauth consent screen
func (o *OAuth) HandleInitCode(code string) (map[string]interface{}, error) {
	if code == "" {
		return nil, fmt.Errorf("invalid code")
	}
	data := url.Values{
		"code":          {code},
		"client_id":     {o.clientID},
		"client_secret": {o.clientSecret},
		"redirect_uri":  {o.redirectURL},
		"grant_type":    {"authorization_code"},
	}
	req, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", data)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
