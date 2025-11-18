package google

import (
	"app/internal/auth"
	"app/internal/user"
	"context"
	"encoding/base64"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

type Service struct {
	oauth *GoogleOAuth	
	jwt *auth.JWTManager
	userRepo *user.Repository
	authRepo *auth.Repository
	dbPool *pgxpool.Pool
}

func NewService(
	oauth *GoogleOAuth, 
	jwt *auth.JWTManager, 
	repo *user.Repository, 
	authRepo *auth.Repository, 
	dbPool *pgxpool.Pool, ) *Service {
	return &Service{oauth, jwt, repo, authRepo, dbPool}
}


func (s *Service) StartGoogleAuth() (state, url string){
	state = auth.GenerateState()
	url = s.oauth.config.AuthCodeURL(state, oauth2.SetAuthURLParam("access_type","offline"))
	return state, url
}

func (s *Service) HandleGoogleCallback(ctx context.Context, code string, ip string, agent string) (auth.GoogleCallbackResp, error ) {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return auth.GoogleCallbackResp{}, fmt.Errorf("failed to begin db pool: %w", err)
	}
	defer tx.Rollback(ctx)
	txAuthRepo := auth.NewRepository(tx)
	txUserRepo := user.NewRepository(tx)

	token, err := s.oauth.ExchangeCode(ctx, code)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	googleUser, err := s.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	u, _ := txUserRepo.GetByEmail(ctx, googleUser.Email)
	if u == nil {
		id, err := txUserRepo.Create(ctx, user.User{
			ID: uuid.New(),
			CreateUser: user.CreateUser{
				Email: googleUser.Email,
				Name: googleUser.Name,
			},
		})
		if err != nil {
		return auth.GoogleCallbackResp{}, err
		}
		u, _ = txUserRepo.GetByID(ctx, id)
	}

	accessToken, err := s.jwt.Generate(u)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	refreshTokenBytes, err := auth.GenerateRandomBytes(32)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}
	if !utf8.ValidString(base64.RawURLEncoding.EncodeToString((auth.HashBytes(refreshTokenBytes)))) {
		fmt.Println("Invalid string hash")
	}

	refreshToken := auth.CreateRefreshToken{
		UserID: u.ID,
		RefreshTokenHash: base64.RawURLEncoding.EncodeToString((auth.HashBytes(refreshTokenBytes))),
		UserAgent: agent,
		// IP: net.ParseIP(ip),
		IP: ip,
		CreatedAt: time.Now(),
		ExpireAt: time.Now().Add(time.Hour * 24 * 30),
	}

	_, err = txAuthRepo.CreateRefToken(ctx, refreshToken)
	if err != nil {
		fmt.Println("CreateRefToken error")
		return auth.GoogleCallbackResp{}, err
	}

	provider := auth.Provider {
		Provider: "google",
		ProviderUserID: googleUser.ID,
		UserID: u.ID,
		CreatedAt: time.Now(),
	}
	_, err = txAuthRepo.CreateAuthProvider(ctx, provider)
	if err != nil {
		fmt.Println("CreateAuthProvider error")
		return auth.GoogleCallbackResp{}, err
	}

	return auth.GoogleCallbackResp{
		User: &user.User{
			CreateUser: user.CreateUser{
				Name: u.Name,
				Email: u.Email,
			},
			ID: u.ID,
		},
		AccessToken: accessToken,
		RefreshToken: string(refreshTokenBytes),
	}, nil
}

