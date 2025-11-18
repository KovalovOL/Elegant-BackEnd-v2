package google

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"app/internal/auth"
	"golang.org/x/oauth2"
	g "golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config  *oauth2.Config
}

func NewGoogleOAuth() (*GoogleOAuth, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectUrl := os.Getenv("REDIRECT_URL")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("client id or secret not found")
	}

	config := &oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		RedirectURL: redirectUrl,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: g.Endpoint,
	}
	return &GoogleOAuth{config}, nil
}

func (g *GoogleOAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

func (g *GoogleOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*auth.UserGoogleResp, error) {
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user auth.UserGoogleResp
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}