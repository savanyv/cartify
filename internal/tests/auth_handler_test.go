package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/delivery/handlers"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type responseBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func setupAuthApp(mockAuth *mocks.AuthUsecaseMock) *fiber.App {
	app := fiber.New()
	authHandler := handlers.NewAuthHandler(mockAuth)
	app.Post("/api/v1/auth/login", authHandler.Login)
	app.Post("/api/v1/auth/register", authHandler.Register)
	app.Post("/api/v1/auth/refresh", authHandler.RefreshToken)

	user := app.Group("/api/v1/user", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})
	user.Get("/profile", authHandler.GetProfile)
	user.Post("/change-password", authHandler.ChangePassword)
	user.Post("/logout", authHandler.Logout)
	return app
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
		expectedTokenKey   string
	}{
		{
			name:        "success",
			requestBody: `{"email":"john@example.com","password":"secret123"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Login", mock.Anything, mock.MatchedBy(func(req dto.LoginRequest) bool {
						return req.Email == "john@example.com" && req.Password == "secret123"
					})).
					Return(&dto.LoginResponse{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
						User: dto.UserResponse{
							ID:       "user-id",
							Name:     "John Doe",
							Username: "johndoe",
							Email:    "john@example.com",
							Role:     "user",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Login successful",
			expectedSuccess:    true,
			expectedTokenKey:   "access_token",
		},
		{
			name:        "validation error",
			requestBody: `{"email":"not-an-email","password":""}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
		{
			name:        "unauthorized invalid credentials",
			requestBody: `{"email":"john@example.com","password":"wrong-password"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Login", mock.Anything, mock.MatchedBy(func(req dto.LoginRequest) bool {
						return req.Email == "john@example.com" && req.Password == "wrong-password"
					})).
					Return(nil, errors.New("invalid credentials"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "invalid credentials",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			if tc.expectedTokenKey != "" {
				data, ok := body.Data.(map[string]interface{})
				require.True(t, ok)
				_, ok = data[tc.expectedTokenKey]
				require.True(t, ok)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			requestBody: `{"name":"John","username":"johndoe","email":"john@example.com","password":"secret123"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Register", mock.Anything, mock.MatchedBy(func(req dto.RegisterRequest) bool {
						return req.Email == "john@example.com" && req.Username == "johndoe"
					})).
					Return(&dto.RegisterResponse{
						Message: "User registered successfully",
						User: dto.UserResponse{
							ID:       "user-id",
							Name:     "John",
							Username: "johndoe",
							Email:    "john@example.com",
							Role:     "user",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "Registration successful",
			expectedSuccess:    true,
		},
		{
			name:        "validation error",
			requestBody: `{"name":"","username":"","email":"bad","password":"12"}`,
			setupMock:   func(mockAuth *mocks.AuthUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
		{
			name:        "email already used",
			requestBody: `{"name":"John","username":"johndoe","email":"used@example.com","password":"secret123"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Register", mock.Anything, mock.MatchedBy(func(req dto.RegisterRequest) bool {
						return req.Email == "used@example.com"
					})).
					Return(nil, errors.New("email already used"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "email already used",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name: "success",
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("GetUserByID", mock.Anything, "test-user-id").
					Return(&dto.UserResponse{
						ID:       "test-user-id",
						Name:     "John",
						Username: "johndoe",
						Email:    "john@example.com",
						Role:     "user",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Profile retrieved",
			expectedSuccess:    true,
		},
		{
			name: "not found",
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("GetUserByID", mock.Anything, "test-user-id").
					Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "user not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			requestBody: `{"old_password":"old123","new_password":"new123"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("ChangePassword", mock.Anything, "test-user-id", mock.MatchedBy(func(req dto.ChangePasswordRequest) bool {
						return req.OldPassword == "old123" && req.NewPassword == "new123"
					})).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Password changed successfully",
			expectedSuccess:    true,
		},
		{
			name:        "invalid old password",
			requestBody: `{"old_password":"wrong","new_password":"new123"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("ChangePassword", mock.Anything, "test-user-id", mock.MatchedBy(func(req dto.ChangePasswordRequest) bool {
						return req.OldPassword == "wrong"
					})).
					Return(errors.New("invalid old password"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid old password",
			expectedSuccess:    false,
		},
		{
			name:        "validation error",
			requestBody: `{"old_password":"","new_password":""}`,
			setupMock:   func(mockAuth *mocks.AuthUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/change-password", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			requestBody: `{"refresh_token":"valid-refresh-token"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("RefreshToken", mock.Anything, "valid-refresh-token").
					Return(&dto.RefreshTokenResponse{
						AccessToken: "new-access-token",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Token refreshed",
			expectedSuccess:    true,
		},
		{
			name:        "invalid token",
			requestBody: `{"refresh_token":"bad-token"}`,
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("RefreshToken", mock.Anything, "bad-token").
					Return(nil, errors.New("invalid refresh token"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "invalid refresh token",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.AuthUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name: "success",
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Logout", mock.Anything, "test-user-id").
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Logout successful",
			expectedSuccess:    true,
		},
		{
			name: "user not found",
			setupMock: func(mockAuth *mocks.AuthUsecaseMock) {
				mockAuth.
					On("Logout", mock.Anything, "test-user-id").
					Return(errors.New("user not found"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "user not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(mocks.AuthUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockAuth)
			}

			app := setupAuthApp(mockAuth)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/logout", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockAuth.AssertExpectations(t)
		})
	}
}
