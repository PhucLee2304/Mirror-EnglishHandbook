package service

import (
	"context"
	"core/config"
	"core/internal/dto"
	"core/internal/model"
	"core/internal/repo"
	"core/pkg/gateways/firebase"
	"core/pkg/gateways/jwt"
	"core/pkg/response"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo       *repo.UserRepository
	firebaseClient *firebase.Client
	cfg            *config.Config
}

func NewAuthService(
	userRepo *repo.UserRepository,
	firebaseClient *firebase.Client,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		firebaseClient: firebaseClient,
		cfg:            cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, body dto.LoginBody) (*dto.LoginResponse, *response.HTTPError) {
	token, err := s.firebaseClient.VerifyIDToken(ctx, body.IDToken)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusUnauthorized,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeInvalidIDToken,
				Detail: err.Error(),
			},
		}
	}

	firebaseUser, err := s.firebaseClient.GetUser(ctx, token.UID)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusUnauthorized,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGetFirebaseUser,
				Detail: err.Error(),
			},
		}
	}

	user, err := s.userRepo.GetByEmail(ctx, firebaseUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var name string
			if body.Name != nil && len(*body.Name) > 0 {
				name = *body.Name
			} else if len(firebaseUser.DisplayName) > 0 {
				name = firebaseUser.DisplayName
			} else {
				name = "Anonymous"
			}

			platform := "firebase"
			if len(firebaseUser.ProviderUserInfo) > 0 {
				platform = firebaseUser.ProviderUserInfo[0].ProviderID
			}

			userModel := &model.User{
				Email:    firebaseUser.Email,
				Name:     name,
				Platform: platform,
			}
			if firebaseUser.PhotoURL != "" {
				photoURL := firebaseUser.PhotoURL
				userModel.Avatar = &photoURL
			}

			createdUser, err := s.userRepo.Create(ctx, userModel)
			if err != nil {
				return nil, &response.HTTPError{
					StatusCode: http.StatusInternalServerError,
					Error: response.ErrorResponse{
						Error:  response.MessageCodeFailedToLogin,
						Detail: err.Error(),
					},
				}
			}
			user = createdUser
		} else {
			return nil, &response.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Error: response.ErrorResponse{
					Error:  response.MessageCodeFailedToLogin,
					Detail: err.Error(),
				},
			}
		}
	}

	accessToken, err := jwt.GenerateAccessToken(s.cfg, user.ID)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGenerateAccessToken,
				Detail: err.Error(),
			},
		}
	}

	refreshToken, err := jwt.GenerateRefreshToken(s.cfg, user.ID)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGenerateRefreshToken,
				Detail: err.Error(),
			},
		}
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, *response.HTTPError) {
	claims, err := jwt.ValidateRefreshToken(s.cfg, refreshToken)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeInvalidRefreshToken,
				Detail: err.Error(),
			},
		}
	}

	user, err := s.userRepo.GetByID(ctx, claims.ID)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeUserNotFound,
				Detail: err.Error(),
			},
		}
	}

	accessToken, err := jwt.GenerateAccessToken(s.cfg, user.ID)
	if err != nil {
		return nil, &response.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Error: response.ErrorResponse{
				Error:  response.MessageCodeFailedToGenerateAccessToken,
				Detail: err.Error(),
			},
		}
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken,
	}, nil
}
