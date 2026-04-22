package usecase

import (
	"cinema-booking-api/internal/auth/dto"
	userDomain "cinema-booking-api/internal/user/domain"
	"cinema-booking-api/pkg/enums"
	"cinema-booking-api/pkg/jwt"
	"cinema-booking-api/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authUsecase struct {
	repo        userDomain.AuthRepository
	jwt         *jwt.JWT
	oauthConfig *oauth2.Config
	redirectURL string
}

func NewAuthUsecase(
	repo userDomain.AuthRepository,
	jwt *jwt.JWT,
	clientID, clientSecret, redirectURL string,
) userDomain.AuthUsecase {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
	}

	return &authUsecase{
		repo:        repo,
		jwt:         jwt,
		oauthConfig: oauthConfig,
		redirectURL: redirectURL,
	}
}

func (u *authUsecase) Register(req interface{}) (interface{}, error) {
	registerReq, ok := req.(dto.RegisterRequest)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	// Check if user already exists
	existingUser, err := u.repo.FindByEmail(registerReq.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(registerReq.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &userDomain.User{
		Username: registerReq.Username,
		Email:    registerReq.Email,
		Password: hashedPassword,
		Role:     enums.UserRoleUser,
		Status:   enums.UserStatusActive,
		Provider: "email",
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := u.jwt.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Provider)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         u.toUserResponse(user),
	}, nil
}

func (u *authUsecase) Login(req interface{}) (interface{}, error) {
	loginReq, ok := req.(dto.LoginRequest)
	if !ok {
		return nil, errors.New("invalid request type")
	}

	user, err := u.repo.FindByEmail(loginReq.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(loginReq.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status != enums.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	// Generate tokens
	tokenPair, err := u.jwt.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Provider)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         u.toUserResponse(user),
	}, nil
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (u *authUsecase) LoginWithGoogle(code string) (interface{}, error) {
	// Exchange code for token
	token, err := u.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	client := u.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Find or create user
	user, err := u.repo.FindByProviderID("google", userInfo.ID)
	if err != nil {
		// User not found, create new one
		user = &userDomain.User{
			Username:   userInfo.Name,
			Email:      userInfo.Email,
			Provider:   "google",
			ProviderID: userInfo.ID,
			AvatarURL:  userInfo.Picture,
			Role:       enums.UserRoleUser,
			Status:     enums.UserStatusActive,
		}

		// Check if email already exists with email provider
		existingUser, _ := u.repo.FindByEmail(userInfo.Email)
		if existingUser != nil {
			// Link to existing account
			existingUser.Provider = "google"
			existingUser.ProviderID = userInfo.ID
			existingUser.AvatarURL = userInfo.Picture
			if err := u.repo.Update(existingUser); err != nil {
				return nil, err
			}
			user = existingUser
		} else {
			if err := u.repo.Create(user); err != nil {
				return nil, err
			}
		}
	}

	// Check if user is active
	if user.Status != enums.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	// Generate tokens
	tokenPair, err := u.jwt.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Provider)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         u.toUserResponse(user),
	}, nil
}

func (u *authUsecase) RefreshToken(refreshToken string) (interface{}, error) {
	claims, err := u.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := u.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != enums.UserStatusActive {
		return nil, errors.New("user account is not active")
	}

	tokenPair, err := u.jwt.GenerateTokenPair(user.ID, user.Email, string(user.Role), user.Provider)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         u.toUserResponse(user),
	}, nil
}

func (u *authUsecase) ValidateToken(token string) (*userDomain.User, error) {
	claims, err := u.jwt.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := u.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *authUsecase) toUserResponse(user *userDomain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		Status:    string(user.Status),
		Provider:  user.Provider,
		AvatarURL: user.AvatarURL,
	}
}

// GetGoogleAuthURL returns the URL for Google OAuth
func (u *authUsecase) GetGoogleAuthURL() string {
	// Generate state for CSRF protection
	state := fmt.Sprintf("%d", time.Now().Unix())
	authURL := u.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return authURL
}
