package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthUsecase() (*usecase.AuthUsecase, *mocks.UserRepositoryMock, *mocks.JWTServiceMock, *mocks.BcryptServiceMock) {
	userRepo := new(mocks.UserRepositoryMock)
	jwtSvc := new(mocks.JWTServiceMock)
	bcryptSvc := new(mocks.BcryptServiceMock)
	return usecase.NewAuthUsecase(userRepo, jwtSvc, bcryptSvc), userRepo, jwtSvc, bcryptSvc
}

func TestAuthUsecase_Register_Success(t *testing.T) {
	uc, userRepo, _, bcryptSvc := setupAuthUsecase()
	ctx := context.Background()

	req := dto.RegisterRequest{
		Name: "Test User", Username: "testuser",
		Email: "test@example.com", Password: "password123",
	}

	userRepo.On("FindByEmail", ctx, req.Email).Return(nil, nil)
	userRepo.On("FindByUsername", ctx, req.Username).Return(nil, nil)
	bcryptSvc.On("HashPassword", req.Password).Return("hashedpass", nil)
	userRepo.On("Create", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.Name == req.Name && u.Username == req.Username &&
			u.Email == req.Email && u.Password == "hashedpass" &&
			u.Role == model.RoleUser
	})).Return(nil)

	resp, err := uc.Register(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", resp.Message)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Email, resp.User.Email)
	userRepo.AssertExpectations(t)
	bcryptSvc.AssertExpectations(t)
}

func TestAuthUsecase_Register_EmailAlreadyUsed(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	req := dto.RegisterRequest{Name: "T", Username: "t", Email: "t@t.com", Password: "pass123"}
	userRepo.On("FindByEmail", ctx, req.Email).Return(&model.User{}, nil)

	resp, err := uc.Register(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "email already used", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_UsernameAlreadyUsed(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	req := dto.RegisterRequest{Name: "T", Username: "t", Email: "t@t.com", Password: "pass123"}
	userRepo.On("FindByEmail", ctx, req.Email).Return(nil, nil)
	userRepo.On("FindByUsername", ctx, req.Username).Return(&model.User{}, nil)

	resp, err := uc.Register(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "username already used", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	uc, userRepo, jwtSvc, bcryptSvc := setupAuthUsecase()
	ctx := context.Background()

	userID := uuid.New()
	req := dto.LoginRequest{Email: "t@t.com", Password: "pass123"}
	user := &model.User{
		ID: userID, Name: "T", Username: "t",
		Email: req.Email, Password: "hashed", Role: model.RoleUser,
	}

	userRepo.On("FindByEmail", ctx, req.Email).Return(user, nil)
	bcryptSvc.On("ComparePassword", user.Password, req.Password).Return(true)
	jwtSvc.On("GenerateAccessToken", userID.String(), user.Username, user.Email, string(user.Role), user.TokenVersion).Return("at", nil)
	jwtSvc.On("GenerateRefreshToken", userID.String()).Return("rt", nil)

	resp, err := uc.Login(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "at", resp.AccessToken)
	assert.Equal(t, "rt", resp.RefreshToken)
	assert.Equal(t, userID.String(), resp.User.ID)
	userRepo.AssertExpectations(t)
	bcryptSvc.AssertExpectations(t)
	jwtSvc.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidCredentials(t *testing.T) {
	uc, userRepo, _, bcryptSvc := setupAuthUsecase()
	ctx := context.Background()

	req := dto.LoginRequest{Email: "t@t.com", Password: "wrong"}
	user := &model.User{ID: uuid.New(), Email: req.Email, Password: "h"}
	userRepo.On("FindByEmail", ctx, req.Email).Return(user, nil)
	bcryptSvc.On("ComparePassword", user.Password, req.Password).Return(false)

	resp, err := uc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
	bcryptSvc.AssertExpectations(t)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	userRepo.On("FindByEmail", ctx, "x@x.com").Return(nil, nil)

	resp, err := uc.Login(ctx, dto.LoginRequest{Email: "x@x.com", Password: "p"})
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserByID_Success(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	uid := uuid.New().String()
	user := &model.User{
		ID: uuid.MustParse(uid), Name: "N", Username: "u",
		Email: "e@e.com", Role: model.RoleUser,
	}

	userRepo.On("FindByID", ctx, uid).Return(user, nil)

	resp, err := uc.GetUserByID(ctx, uid)
	assert.NoError(t, err)
	assert.Equal(t, uid, resp.ID)
	assert.Equal(t, "N", resp.Name)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserByID_NotFound(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	userRepo.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.GetUserByID(ctx, "x")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, resp)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_ChangePassword_Success(t *testing.T) {
	uc, userRepo, _, bcryptSvc := setupAuthUsecase()
	ctx := context.Background()

	uid := uuid.New().String()
	req := dto.ChangePasswordRequest{OldPassword: "old", NewPassword: "new"}
	userRepo.On("FindByID", ctx, uid).Return(&model.User{ID: uuid.MustParse(uid), Password: "h"}, nil)
	bcryptSvc.On("ComparePassword", "h", "old").Return(true)
	bcryptSvc.On("HashPassword", "new").Return("hn", nil)
	userRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
		return u.Password == "hn"
	})).Return(nil)

	err := uc.ChangePassword(ctx, uid, req)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	bcryptSvc.AssertExpectations(t)
}

func TestAuthUsecase_ChangePassword_InvalidOldPassword(t *testing.T) {
	uc, userRepo, _, bcryptSvc := setupAuthUsecase()
	ctx := context.Background()

	uid := uuid.New().String()
	userRepo.On("FindByID", ctx, uid).Return(&model.User{ID: uuid.MustParse(uid), Password: "h"}, nil)
	bcryptSvc.On("ComparePassword", "h", "wrong").Return(false)

	err := uc.ChangePassword(ctx, uid, dto.ChangePasswordRequest{OldPassword: "wrong", NewPassword: "n"})
	assert.Error(t, err)
	assert.Equal(t, "invalid old password", err.Error())
	userRepo.AssertExpectations(t)
	bcryptSvc.AssertExpectations(t)
}

func TestAuthUsecase_ChangePassword_UserNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	userRepo.On("FindByID", ctx, "x").Return(nil, nil)

	err := uc.ChangePassword(ctx, "x", dto.ChangePasswordRequest{})
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_RefreshToken_Success(t *testing.T) {
	uc, userRepo, jwtSvc, _ := setupAuthUsecase()
	ctx := context.Background()

	userID := uuid.New()
	rc := jwt.RegisteredClaims{Subject: userID.String()}

	jwtSvc.On("ValidateRefreshToken", "rt").Return(&rc, nil)
	userRepo.On("FindByID", ctx, userID.String()).Return(
		&model.User{ID: userID, Username: "u", Email: "e@e.com", Role: model.RoleUser}, nil,
	)
	jwtSvc.On("GenerateAccessToken", userID.String(), "u", "e@e.com", string(model.RoleUser), 0).Return("new_at", nil)

	resp, err := uc.RefreshToken(ctx, "rt")
	assert.NoError(t, err)
	assert.Equal(t, "new_at", resp.AccessToken)
	userRepo.AssertExpectations(t)
	jwtSvc.AssertExpectations(t)
}

func TestAuthUsecase_RefreshToken_InvalidToken(t *testing.T) {
	uc, _, jwtSvc, _ := setupAuthUsecase()
	ctx := context.Background()

	jwtSvc.On("ValidateRefreshToken", "bad").Return(nil, errors.New("invalid"))

	resp, err := uc.RefreshToken(ctx, "bad")
	assert.Error(t, err)
	assert.Equal(t, "invalid refresh token", err.Error())
	assert.Nil(t, resp)
	jwtSvc.AssertExpectations(t)
}

func TestAuthUsecase_Logout_Success(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	uid := uuid.New().String()
	userRepo.On("FindByID", ctx, uid).Return(&model.User{ID: uuid.MustParse(uid), TokenVersion: 0}, nil)
	userRepo.On("UpdateTokenVersion", ctx, uid, 1).Return(nil)

	err := uc.Logout(ctx, uid)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthUsecase_Logout_UserNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase()
	ctx := context.Background()

	userRepo.On("FindByID", ctx, "x").Return(nil, nil)

	err := uc.Logout(ctx, "x")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	userRepo.AssertExpectations(t)
}
