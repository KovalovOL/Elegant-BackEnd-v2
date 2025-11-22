package google

import (
	"app/internal/auth"
	"app/internal/user"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	userRepo *user.Repository, 
	authRepo *auth.Repository, 
	dbPool *pgxpool.Pool, ) *Service {
	return &Service{oauth, jwt, userRepo, authRepo, dbPool}
}


func (s *Service) StartGoogleAuth() (state, url string){
	state = auth.GenerateState()
	url = s.oauth.config.AuthCodeURL(state, oauth2.SetAuthURLParam("access_type", "offline"))
	return state, url
}

func (s *Service) getGoogleUser(ctx context.Context, code string) (*auth.UserGoogleResp, error) {
	token, err := s.oauth.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}
	googleUser, err := s.oauth.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	return googleUser, nil
}

func (s *Service) getAndCreateUser(
	ctx context.Context, 
	txUserRepo user.Repository, 
	googleUser auth.UserGoogleResp,
) (*user.User, error) {

	u, _ := txUserRepo.GetByEmail(ctx, googleUser.Email)
	if u == nil {
		u = &user.User{
			ID: uuid.New(),
			CreateUser: user.CreateUser{
				Email: googleUser.Email,
				Name: googleUser.Name,
			},
		}
		_, err := txUserRepo.Create(ctx, *u)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (s *Service) saveRefreshToken(
	ctx context.Context, 
	refreshTokenBytes []byte,
	txAuthRepo auth.Repository, 
	userID uuid.UUID,
	googleID string,
	userAgent string,
	userIP net.IP,
) error {

	refreshTokenHash := base64.RawURLEncoding.EncodeToString((auth.HashBytes(refreshTokenBytes)))

	refreshToken := auth.CreateRefreshToken{
		UserID: userID,
		RefreshTokenHash: refreshTokenHash,
		UserAgent: userAgent,
		IP: userIP,
		CreatedAt: time.Now().UTC(),
		ExpireAt: time.Now().UTC().Add(time.Hour * 24 * 30),
	}
	_, err := txAuthRepo.CreateRefToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	provider := auth.Provider {
		Provider: "google",
		ProviderUserID: googleID,
		UserID: userID,
		CreatedAt: time.Now().UTC(),
	}
	_, err = txAuthRepo.CreateAuthProvider(ctx, provider)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) isRelogin(ctx context.Context, refTCookieBytes []byte, userID uuid.UUID) bool {
	refTokenCookieHash := auth.HashBytes(refTCookieBytes)
	oldRefToken, err := s.authRepo.GetRefToken(ctx, base64.RawURLEncoding.EncodeToString(refTokenCookieHash))
	if err == nil {
		if oldRefToken.UserID == userID {
			return true
		}
	}
	return false
}

func (s *Service) HandleGoogleCallback(
	ctx context.Context,
	refTokenCookie string,
	code string, 
	ip string, 
	agent string,
) (auth.GoogleCallbackResp, error ) {

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return auth.GoogleCallbackResp{}, fmt.Errorf("failed to begin db pool: %w", err)
	}
	defer func() {
    	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
        	log.Printf("rollback error: %v", err)
    	}
	}()

	txAuthRepo := auth.NewRepository(tx)
	txUserRepo := user.NewRepository(tx)

	googleUser, err := s.getGoogleUser(ctx, code)
	if err != nil {
		return auth.GoogleCallbackResp{}, err		
	}

	u, err := s.getAndCreateUser(ctx, *txUserRepo, *googleUser)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	accToken, err := s.jwt.Generate(u)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}
	refTokenBytes, err := auth.GenerateRandomBytes(32)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	refTokenCookieBytes, err :=  base64.RawURLEncoding.DecodeString(refTokenCookie)
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}
	if s.isRelogin(ctx, refTokenCookieBytes, u.ID) {
		refTokenBytes = refTokenCookieBytes
	}

	err = s.saveRefreshToken(ctx, refTokenBytes, *txAuthRepo, u.ID, googleUser.ID, agent, net.ParseIP(ip))
	if err != nil {
		return auth.GoogleCallbackResp{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return auth.GoogleCallbackResp{}, err
	}
	return auth.GoogleCallbackResp{
		User: u,
		AccessToken: accToken,
		RefreshToken: base64.RawURLEncoding.EncodeToString(refTokenBytes),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refToken string) (string, error) {
	refTokenBytes, err :=  base64.RawURLEncoding.DecodeString(refToken)
	if err != nil {
		return "", fmt.Errorf("failed to decode refresh token: %w", err)
	}

	refTokenHash := auth.HashBytes(refTokenBytes)
	refTokenHashStr := base64.RawURLEncoding.EncodeToString(refTokenHash)
	t, err := s.authRepo.GetRefToken(ctx, refTokenHashStr)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	if time.Now().UTC().After(t.ExpireAt) {
		return "", fmt.Errorf("refresh token is expired")
	}

	user, err := s.userRepo.GetByID(ctx, t.UserID)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	
	accToken, err := s.jwt.Generate(user)
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	if err = s.authRepo.RefreshTokenTime(ctx, refTokenHashStr); err != nil {
		return "", fmt.Errorf("failed to update refresh token: %w", err)
	}

	return accToken, err
}

func (s *Service) Logout(ctx context.Context, refTokenCookie string) error {
	refTokenBytes, err :=  base64.RawURLEncoding.DecodeString(refTokenCookie)
	if err != nil {
		return err
	}

	refTokenHash := auth.HashBytes(refTokenBytes)
	refTokenHashStr := base64.RawURLEncoding.EncodeToString(refTokenHash)
	err = s.authRepo.DeleteRefTokenByHash(ctx, refTokenHashStr)
	if err != nil {
		return err
	}
	return nil
}
