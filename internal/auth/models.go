package auth

import (
	"app/internal/user"
	"time"
	"net"
	"github.com/google/uuid"
)


type UserGoogleResp struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleCallbackResp struct {
	User 		 *user.User `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateRefreshToken struct {
	UserID 			 uuid.UUID
	RefreshTokenHash string
	UserAgent 		 string
	IP 				 net.IP
	ExpireAt 		 time.Time
	CreatedAt 		 time.Time
}

type Provider struct {
	Provider 	   string
	ProviderUserID string //Id recieved from provider
	UserID 		   uuid.UUID
	CreatedAt 	   time.Time
}