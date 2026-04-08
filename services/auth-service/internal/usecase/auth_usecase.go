package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/afifn11/gopay-x/services/auth-service/config"
	"github.com/afifn11/gopay-x/services/auth-service/internal/domain"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *LoginRequest) (*domain.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	ValidateToken(ctx context.Context, accessToken string) (*domain.Claims, error)
}

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type authUsecase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	cfg       *config.Config
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	cfg *config.Config,
) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *RegisterRequest) (*domain.User, error) {
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:       uuid.New(),
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *LoginRequest) (*domain.TokenPair, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	_ = uc.userRepo.UpdateLastLogin(ctx, user.ID)

	return tokenPair, nil
}

func (uc *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	storedToken, err := uc.tokenRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil || storedToken == nil {
		return nil, ErrInvalidToken
	}

	user, err := uc.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if err := uc.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return uc.generateTokenPair(ctx, user)
}

func (uc *authUsecase) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if err := uc.tokenRepo.BlacklistAccessToken(ctx, accessToken, uc.cfg.JWT.AccessExpiry); err != nil {
		return err
	}
	return uc.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (uc *authUsecase) ValidateToken(ctx context.Context, accessToken string) (*domain.Claims, error) {
	blacklisted, err := uc.tokenRepo.IsBlacklisted(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, ErrInvalidToken
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.cfg.JWT.AccessSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &domain.Claims{
		UserID: mapClaims["user_id"].(string),
		Email:  mapClaims["email"].(string),
		Role:   domain.UserRole(mapClaims["role"].(string)),
	}, nil
}

func (uc *authUsecase) generateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	// Access token
	accessExpiry := time.Duration(uc.cfg.JWT.AccessExpiry) * time.Minute
	accessClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
		"exp":     time.Now().Add(accessExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(uc.cfg.JWT.AccessSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshExpiry := time.Duration(uc.cfg.JWT.RefreshExpiry) * 24 * time.Hour
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(refreshExpiry).Unix(),
	}
	refreshTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(uc.cfg.JWT.RefreshSecret))
	if err != nil {
		return nil, err
	}

	// Store refresh token in DB
	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(refreshExpiry),
	}
	if err := uc.tokenRepo.StoreRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    uc.cfg.JWT.AccessExpiry * 60,
	}, nil
}