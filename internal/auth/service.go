package auth

import (
	"app/internal/user"
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Service struct {
	oauth *GoogleOAuth	
	jwt *JWTManager
	userRepo *user.Repository
}

func NewService(oauth *GoogleOAuth, jwt *JWTManager, repo *user.Repository) *Service {
	return &Service{oauth, jwt, repo}
}


func (s *Service) StartGoogleAuth() (state, url string){
	state = generateState()
	url = s.oauth.config.AuthCodeURL(state, oauth2.SetAuthURLParam("access_type","offline"))
	return state, url
}

func (s *Service) HandleGoogleCallback(ctx context.Context, code string) (*user.User, string, error) {
	token, err := s.oauth.ExchangeCode(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code")
	}

	googleUser, err := s.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, "", err
	}

	u, _ := s.userRepo.GetByEmail(ctx, googleUser.Email)
	if u == nil {
		id, err := s.userRepo.Create(ctx, user.User{
			ID: uuid.New(),
			CreateUser: user.CreateUser{
				Email: googleUser.Email,
				Name: googleUser.Name,
			},
		})
		if err != nil {
			return nil, "", err
		}
		u, _ = s.userRepo.GetByID(ctx, id)
	}

	jwtToken, err := s.jwt.Generate(u)
	if err != nil {
		return nil, "", err
	}

	return u, jwtToken, nil	
}

