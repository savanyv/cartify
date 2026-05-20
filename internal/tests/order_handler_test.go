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
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupOrderApp(mockOrder *mocks.OrderUsecaseMock) *fiber.App {
	app := fiber.New()
	orderHandler := handlers.NewOrderHandler(mockOrder)

	orders := app.Group("/api/v1/orders", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/", orderHandler.GetUserOrders)
	orders.Get("/:id", orderHandler.GetOrderByID)

	admin := app.Group("/api/v1/admin", func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})
	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Put("/orders/:id/status", orderHandler.UpdateOrderStatus)
	return app
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(*mocks.OrderUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name: "success",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("CreateOrder", mock.Anything, "test-user-id").
					Return(&dto.OrderResponse{
						ID:         "order-123",
						Status:     "pending",
						TotalPrice: 50000,
						Items: []dto.OrderItemResponse{
							{ID: "oi-1", Qty: 2, Price: 25000, SubTotal: 50000},
						},
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "Order created successfully",
			expectedSuccess:    true,
		},
		{
			name: "empty cart",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("CreateOrder", mock.Anything, "test-user-id").
					Return(nil, errors.New("cart is empty"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "cart is empty",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrder := new(mocks.OrderUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockOrder)
			}

			app := setupOrderApp(mockOrder)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockOrder.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_GetUserOrders(t *testing.T) {
	tests := []struct {
		name               string
		queryString        string
		setupMock          func(*mocks.OrderUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			queryString: "",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("GetuserOrders", mock.Anything, "test-user-id", 1, 10, "", "created_at", "desc").
					Return([]dto.OrderResponse{
						{ID: "order-1", Status: "pending", TotalPrice: 25000},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Orders retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name:        "empty list",
			queryString: "",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("GetuserOrders", mock.Anything, "test-user-id", 1, 10, "", "created_at", "desc").
					Return([]dto.OrderResponse{}, int64(0), nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Orders retrieved successfully",
			expectedSuccess:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrder := new(mocks.OrderUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockOrder)
			}

			app := setupOrderApp(mockOrder)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders"+tc.queryString, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockOrder.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_GetOrderByID(t *testing.T) {
	tests := []struct {
		name               string
		orderID            string
		setupMock          func(*mocks.OrderUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:    "success",
			orderID: "order-123",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("GetOrderByID", mock.Anything, "test-user-id", "order-123").
					Return(&dto.OrderResponse{
						ID:         "order-123",
						Status:     "pending",
						TotalPrice: 50000,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Order retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name:    "not found",
			orderID: "missing-order",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("GetOrderByID", mock.Anything, "test-user-id", "missing-order").
					Return(nil, errors.New("order not found"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "order not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrder := new(mocks.OrderUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockOrder)
			}

			app := setupOrderApp(mockOrder)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/"+tc.orderID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockOrder.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_GetAllOrders(t *testing.T) {
	tests := []struct {
		name               string
		queryString        string
		setupMock          func(*mocks.OrderUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			queryString: "",
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("GetAllOrders", mock.Anything, 1, 10, "", "created_at", "desc").
					Return([]dto.AdminOrderResponse{
						{
							ID:     "order-1",
							Status: "pending",
							UserID: "user-1",
						},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "All orders retrieved successfully",
			expectedSuccess:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrder := new(mocks.OrderUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockOrder)
			}

			app := setupOrderApp(mockOrder)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/orders"+tc.queryString, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body responseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockOrder.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
	tests := []struct {
		name               string
		orderID            string
		requestBody        string
		setupMock          func(*mocks.OrderUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			orderID:     "order-123",
			requestBody: `{"status":"paid"}`,
			setupMock: func(mockOrder *mocks.OrderUsecaseMock) {
				mockOrder.
					On("UpdateOrderStatus", mock.Anything, "order-123", model.OrderStatus("paid")).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Order status updated successfully",
			expectedSuccess:    true,
		},
		{
			name:        "invalid status",
			orderID:     "order-123",
			requestBody: `{"status":"invalid_status"}`,
			setupMock:   func(mockOrder *mocks.OrderUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrder := new(mocks.OrderUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockOrder)
			}

			app := setupOrderApp(mockOrder)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/orders/"+tc.orderID+"/status", strings.NewReader(tc.requestBody))
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

			mockOrder.AssertExpectations(t)
		})
	}
}
