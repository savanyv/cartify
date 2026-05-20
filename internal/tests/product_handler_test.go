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

type productResponseBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func setupProductApp(mockProduct *mocks.ProductUsecaseMock) *fiber.App {
	app := fiber.New()
	productHandler := handlers.NewProductHandler(mockProduct)
	app.Get("/api/v1/products/:id", productHandler.GetProductByID)
	app.Get("/api/v1/products", productHandler.GetAllProducts)
	app.Post("/api/v1/admin/products", productHandler.CreateProduct)
	app.Put("/api/v1/admin/products/:id", productHandler.UpdateProduct)
	app.Delete("/api/v1/admin/products/:id", productHandler.DeleteProduct)
	app.Post("/api/v1/admin/products/:product_id/variants", productHandler.CreateVariant)
	app.Put("/api/v1/admin/products/variants/:id", productHandler.UpdateVariant)
	return app
}

func TestProductHandler_GetProductByID(t *testing.T) {
	tests := []struct {
		name               string
		productID          string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:      "success",
			productID: "prod-123",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("GetProductByID", mock.Anything, "prod-123").
					Return(&dto.ProductResponse{
						ID:          "prod-123",
						Name:        "Test Product",
						Description: "A product for testing",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Product retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name:      "not found",
			productID: "missing-id",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("GetProductByID", mock.Anything, "missing-id").
					Return(nil, errors.New("product not found"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "product not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+tc.productID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	tests := []struct {
		name               string
		queryString        string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			queryString: "",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("GetAllProductsWithPagination", mock.Anything, 1, 10, "", "created_at", "desc").
					Return([]dto.ProductResponse{
						{ID: "prod-1", Name: "Product 1"},
						{ID: "prod-2", Name: "Product 2"},
					}, int64(2), nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Products retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name:        "with search and pagination",
			queryString: "?search=phone&page=1&limit=5&sort=name&order=asc",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("GetAllProductsWithPagination", mock.Anything, 1, 5, "phone", "name", "asc").
					Return([]dto.ProductResponse{
						{ID: "prod-3", Name: "Phone"},
					}, int64(1), nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Products retrieved successfully",
			expectedSuccess:    true,
		},
		{
			name:        "internal server error",
			queryString: "",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("GetAllProductsWithPagination", mock.Anything, 1, 10, "", "created_at", "desc").
					Return(nil, int64(0), errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "database error",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/products"+tc.queryString, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_CreateProduct(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			requestBody: `{"name":"New Product","description":"A new product"}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("CreateProduct", mock.Anything, mock.MatchedBy(func(req dto.CreateProductRequest) bool {
						return req.Name == "New Product"
					})).
					Return(&dto.ProductResponse{
						ID:          "new-prod-id",
						Name:        "New Product",
						Description: "A new product",
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "Product created successfully",
			expectedSuccess:    true,
		},
		{
			name:        "validation error",
			requestBody: `{"name":"ab"}`,
			setupMock:   func(mockProduct *mocks.ProductUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
		{
			name:        "bad request",
			requestBody: `{"name":"","description":""}`,
			setupMock:   func(mockProduct *mocks.ProductUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	tests := []struct {
		name               string
		productID          string
		requestBody        string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			productID:   "prod-123",
			requestBody: `{"name":"Updated Name"}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("UpdateProduct", mock.Anything, "prod-123", mock.MatchedBy(func(req dto.UpdateProductRequest) bool {
						return req.Name == "Updated Name"
					})).
					Return(&dto.ProductResponse{
						ID:   "prod-123",
						Name: "Updated Name",
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Product updated successfully",
			expectedSuccess:    true,
		},
		{
			name:        "not found",
			productID:   "missing-id",
			requestBody: `{"name":"Updated"}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("UpdateProduct", mock.Anything, "missing-id", mock.Anything).
					Return(nil, errors.New("product not found"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "product not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/"+tc.productID, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	tests := []struct {
		name               string
		productID          string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:      "success",
			productID: "prod-123",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("DeleteProduct", mock.Anything, "prod-123").
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Product deleted successfully",
			expectedSuccess:    true,
		},
		{
			name:      "not found",
			productID: "missing-id",
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("DeleteProduct", mock.Anything, "missing-id").
					Return(errors.New("product not found"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "product not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/products/"+tc.productID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_CreateVariant(t *testing.T) {
	tests := []struct {
		name               string
		productID          string
		requestBody        string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			productID:   "prod-123",
			requestBody: `{"name":"32GB","stock":10,"price":150000}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("CreateVariant", mock.Anything, "prod-123", mock.MatchedBy(func(req dto.CreateVariantRequest) bool {
						return req.Name == "32GB" && req.Price == 150000
					})).
					Return(&dto.VariantResponse{
						ID:        "var-456",
						Name:      "32GB",
						Stock:     10,
						Price:     150000,
						ProductID: "prod-123",
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "Variant created successfully",
			expectedSuccess:    true,
		},
		{
			name:        "validation error",
			productID:   "prod-123",
			requestBody: `{"name":"","price":0}`,
			setupMock:   func(mockProduct *mocks.ProductUsecaseMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Validation failed",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/products/"+tc.productID+"/variants", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}

func TestProductHandler_UpdateVariant(t *testing.T) {
	tests := []struct {
		name               string
		variantID          string
		requestBody        string
		setupMock          func(*mocks.ProductUsecaseMock)
		expectedStatusCode int
		expectedMessage    string
		expectedSuccess    bool
	}{
		{
			name:        "success",
			variantID:  "var-456",
			requestBody: `{"name":"64GB","price":200000}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("UpdateVariant", mock.Anything, "var-456", mock.MatchedBy(func(req dto.UpdateVariantRequest) bool {
						return req.Name == "64GB" && req.Price == 200000
					})).
					Return(&dto.VariantResponse{
						ID:    "var-456",
						Name:  "64GB",
						Price: 200000,
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Variant updated successfully",
			expectedSuccess:    true,
		},
		{
			name:        "not found",
			variantID:  "missing-var",
			requestBody: `{"name":"64GB"}`,
			setupMock: func(mockProduct *mocks.ProductUsecaseMock) {
				mockProduct.
					On("UpdateVariant", mock.Anything, "missing-var", mock.Anything).
					Return(nil, errors.New("variant not found"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "variant not found",
			expectedSuccess:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockProduct := new(mocks.ProductUsecaseMock)
			if tc.setupMock != nil {
				tc.setupMock(mockProduct)
			}

			app := setupProductApp(mockProduct)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/products/variants/"+tc.variantID, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var body productResponseBody
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSuccess, body.Success)
			require.Equal(t, tc.expectedMessage, body.Message)

			mockProduct.AssertExpectations(t)
		})
	}
}
