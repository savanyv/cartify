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

func setupCartApp(mockCart *mocks.CartUsecaseMock) *fiber.App {
	app := fiber.New()
	cartHandler := handlers.NewCartHandler(mockCart)

	cart := app.Group("/api/v1/cart", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/", cartHandler.AddToCart)
	cart.Put("/items/:item_id", cartHandler.UpdateCartItem)
	cart.Delete("/items/:item_id", cartHandler.RemoveFromCart)
	cart.Delete("/clear", cartHandler.ClearCart)
	return app
}

func TestCartHandler_GetCart(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.CartUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name: "success with items",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("GetCart", mock.Anything, "test-user-id").
					Return(&dto.CartResponse{
						ID:         "cart-123",
						TotalPrice: 50000,
						ItemCount:  2,
						Items: []dto.CartItemResponse{
							{ID: "item-1", Quantity: 1, Price: 25000, SubTotal: 25000},
							{ID: "item-2", Quantity: 1, Price: 25000, SubTotal: 25000},
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Cart retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name: "empty cart",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("GetCart", mock.Anything, "test-user-id").
					Return(&dto.CartResponse{
						ID:         "cart-123",
						TotalPrice: 0,
						ItemCount:  0,
						Items:      []dto.CartItemResponse{},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Cart retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name: "internal server error",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("GetCart", mock.Anything, "test-user-id").
					Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "database error",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCart := new(mocks.CartUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockCart)
			}

			app := setupCartApp(mockCart)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockCart.AssertExpectations(t)
		})
	}
}

func TestCartHandler_AddToCart(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.CartUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			requestBody: `{"product_variant_id":"var-456","quantity":2}`,
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("AddToCart", mock.Anything, "test-user-id", mock.MatchedBy(func(req dto.AddToCartRequest) bool {
						return req.ProductVariantID == "var-456" && req.Quantity == 2
					})).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Item added to cart successfully",
			expectedSuccess:    true,
		},
		{
			name:        "insufficient stock",
			requestBody: `{"product_variant_id":"var-456","quantity":999}`,
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("AddToCart", mock.Anything, "test-user-id", mock.Anything).
					Return(errors.New("insufficient stock"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "insufficient stock",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCart := new(mocks.CartUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockCart)
			}

			app := setupCartApp(mockCart)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/cart", strings.NewReader(tc.requestBody))
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

			mockCart.AssertExpectations(t)
		})
	}
}

func TestCartHandler_UpdateCartItem(t *testing.T) {
	tests := []struct {
		name               string
		itemID             string
		requestBody        string
		setupMock          func(*mocks.CartUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			itemID:      "item-1",
			requestBody: `{"quantity":3}`,
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("UpdateCartItem", mock.Anything, "test-user-id", "item-1", mock.MatchedBy(func(req dto.UpdateCartItemRequest) bool {
						return req.Quantity == 3
					})).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Cart item updated successfully",
			expectedSuccess:    true,
		},
		{
			name:        "item not found",
			itemID:      "missing-item",
			requestBody: `{"quantity":3}`,
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("UpdateCartItem", mock.Anything, "test-user-id", "missing-item", mock.Anything).
					Return(errors.New("item not found in cart"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "item not found in cart",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCart := new(mocks.CartUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockCart)
			}

			app := setupCartApp(mockCart)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/"+tc.itemID, strings.NewReader(tc.requestBody))
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

			mockCart.AssertExpectations(t)
		})
	}
}

func TestCartHandler_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name               string
		itemID             string
		setupMock          func(*mocks.CartUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:   "success",
			itemID: "item-1",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("RemoveFromCart", mock.Anything, "test-user-id", "item-1").
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Item removed from cart successfully",
			expectedSuccess:    true,
		},
		{
			name:   "item not found",
			itemID: "missing-item",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("RemoveFromCart", mock.Anything, "test-user-id", "missing-item").
					Return(errors.New("item not found in cart"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "item not found in cart",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCart := new(mocks.CartUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockCart)
			}

			app := setupCartApp(mockCart)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/"+tc.itemID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockCart.AssertExpectations(t)
		})
	}
}

func TestCartHandler_ClearCart(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.CartUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name: "success",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("ClearCart", mock.Anything, "test-user-id").
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Cart cleared successfully",
			expectedSuccess:    true,
		},
		{
			name: "internal server error",
			setupMock: func(mockCart *mocks.CartUsecaseMock) {
				mockCart.
					On("ClearCart", mock.Anything, "test-user-id").
					Return(errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "database error",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCart := new(mocks.CartUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockCart)
			}

			app := setupCartApp(mockCart)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/clear", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockCart.AssertExpectations(t)
		})
	}
}
