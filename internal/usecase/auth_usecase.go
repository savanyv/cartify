package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type AuthUsecase struct {
	userRepo model.UserRepository
	jwtService helpers.JWTService
	bcryptService helpers.BcryptService
}

func NewAuthUsecase(ur model.UserRepository, js helpers.JWTService, bs helpers.BcryptService) *AuthUsecase {
	return &AuthUsecase{
		userRepo: ur,
		jwtService: js,
		bcryptService: bs,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already used")
	}

	existingUser, err = u.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already used")
	}

	hashedPassword, err := u.bcryptService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: hashedPassword,
		Role: model.RoleUser,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	res := &dto.RegisterResponse{
		Message: "User registered successfully",
		User: dto.UserResponse{
			ID: user.ID.String(),
			Name: user.Name,
			Username: user.Username,
			Email: user.Email,
			Role: string(user.Role),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	return res, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !u.bcryptService.ComparePassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := u.jwtService.GenerateAccessToken(
		user.ID.String(),
		user.Username,
		user.Email,
		string(user.Role),
		user.TokenVersion,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	res := &dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID: user.ID.String(),
			Name: user.Name,
			Username: user.Username,
			Email: user.Email,
			Role: string(user.Role),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}

	return res, nil
}

func (u *AuthUsecase) GetUserByID(ctx context.Context, ID string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	res := &dto.UserResponse{
		ID: user.ID.String(),
		Name: user.Name,
		Username: user.Username,
		Email: user.Email,
		Role: string(user.Role),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	return res, nil
}

func (u *AuthUsecase) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if !u.bcryptService.ComparePassword(user.Password, req.OldPassword) {
		return errors.New("invalid old password")
	}

	newHashedPassword, err := u.bcryptService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = newHashedPassword
	return u.userRepo.Update(ctx, user)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	claims, err := u.jwtService.ValidateRefreshToken(refreshToken)	
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	user, err := u.userRepo.FindByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := u.jwtService.GenerateAccessToken(
		user.ID.String(),
		user.Username,
		user.Email,
		string(user.Role),
		user.TokenVersion,
	)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, userID string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	newVersion := user.TokenVersion + 1
	return u.userRepo.UpdateTokenVersion(ctx, userID, newVersion)
}