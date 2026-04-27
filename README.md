# 🛍️ Cartify - E-commerce Backend API

Cartify adalah backend e-commerce yang dibangun dengan **Go**, **Fiber**, **GORM**, dan **PostgreSQL**.

---

## ✨ Fitur

### Core Features
- **Authentication & Authorization**
  - Register, Login, Logout
  - JWT Access Token (24 jam) & Refresh Token (7 hari)
  - Role-based access (Admin & User)
  - Token version system for secure logout

- **Product Management**
  - CRUD operations (Admin only)
  - Product variants (size, color, stock, price)
  - Pagination, search, sort
  - Public access for viewing products

- **Shopping Cart**
  - Add, update, remove items
  - Clear cart
  - Stock validation
  - Price stored at add time

- **Order Management**
  - Checkout from cart
  - Order history with pagination
  - Order status management (Admin)
  - Stock reduction on order
  - Auto-clear cart after checkout

### Security Features
- Security headers (CSP, HSTS, XSS Protection)
- Request ID for tracing
- Rate limiting
- CORS configuration
- API Key authentication

### Observability
- Structured logging
- Request ID tracking
- Error logging with stack trace

---

## 🛠 Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.21+ |
| **Framework** | Fiber v2 |
| **Database** | PostgreSQL |
| **ORM** | GORM |
| **Auth** | JWT, bcrypt |
| **Config** | godotenv |

---

## 🚀 Installation

### 1. Clone Repository

```bash
git clone https://github.com/savanyv/cartify.git
cd cartify
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Environment

```bash
cp .env.sample .env
#Edit .env with your configuration
```

### 4. Run the Application

```bash
go run cmd/main.go
```
Server akan berjalan di ```http://localhost:8000```

## ⚙️ Environment Variables

```env
# App
APP_NAME=Cartify
APP_ENV=development
APP_PORT=8000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cartify

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# API Key
API_KEY=cartify-frontend-key-2024
```

## 📚 API Documentation

### Base URL

```text
http://localhost:8000
```

### Authentication
##### All Protected endpoints require :
- ```X-API-Key``` header
- ```Authorization: Bearer <access_token>``` header (after login)

### Public Endpoints

| METHOD | Endpoint | Description |
|--------|----------|-------------|
| **GET** | `/health` | Health check |
| **POST** | `/api/v1/auth/register` | Register new user |
| **POST** | `/api/v1/auth/login` | Login user |
| **POST** | `/api/v1/auth/refresh` | Refresh access token |
| **GET** | `/api/v1/products` | Get all products (paginated) |
| **GET** | `/api/v1/products/:id` | Get product by ID |

### Protected Endpoints (User)

| METHOD | Endpoint | Description |
|--------|----------|-------------|
| **GET** | `/api/v1/user/profile` | Get user profile |
| **POST** | `/api/v1/user/change-password` | Change password |
| **POST** | `/api/v1/user/logout` | Logout user |
| **GET** | `/api/v1/cart` | Get user cart |
| **POST** | `/api/v1/cart` | Add item to cart |
| **PUT** | `/api/v1/cart/items/:item_id` | Update cart item |
| **DELETE** | `/api/v1/cart/items/:item_id` | Remove cart item |
| **DELETE** | `/api/v1/cart/clear` | Clear cart |
| **POST** | `/api/v1/orders` | Create order |
| **GET** | `/api/v1/orders` | Get user orders |
| **GET** | `/api/v1/orders/:id` | Get order detail |

### Admin Endpoints

| METHOD | Endpoint | Description |
|--------|----------|-------------|
| **POST** | `/api/v1/admin/products` | Create product |
| **PUT** | `/api/v1/admin/products/:id` | Update product |
| **DELETE** | `/api/v1/admin/products/:id` | Delete product |
| **POST** | `/api/v1/admin/products/:product_id/variants` | Create variant |
| **PUT** | `/api/v1/admin/products/variants/:id` | Update variant |
| **GET** | `/api/v1/admin/orders` | Get all orders |
| **PUT** | `/api/v1/admin/orders/:id/status` | Update order status |

### Query Parameters (Pagination)

| Parameters | Default | Description |
|----------|------------|---------|
| ```page``` | 1 | Page number |
| ```limit``` | 10 | items per page (Max 100) |
| ```search``` | "" | Search keyword |
| ```sort``` | ```created_at``` | Sort field |
| ```order``` | ```desc``` | Sort order (asc/desc) |

### 📁 Project Structure

```text
cartify/
├── cmd/api/main.go                 # Entry point
├── internal/
│   ├── config/                     # Configuration
│   ├── delivery/
│   │   ├── handlers/               # HTTP handlers
│   │   └── routes/                 # Routes
│   ├── infrastructure/             # Database, seeders
│   ├── middlewares/                # HTTP middlewares
│   ├── model/                      # Database models
│   ├── repository/                 # Data access layer
│   ├── usecase/                    # Business logic
│   └── utils/helpers/              # Helper functions
├── .env                            # Environment variables
├── go.mod
└── README.md
```

### 🔒 Default Admin Account (Development)

After first run, admin account is auto-created:

| Field | Value |
|-------|-------|
| Email | superadmin@cartify.com |
| Password | superadmin123! |

----

### 📝 Response Format

#### Success Response (Non-paginated)

```json
{
  "success": true,
  "message": "Login successful",
  "data": { ... }
}
```

#### Success Response (Paginated)

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [ ... ],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "total_pages": 10,
    "has_prev": false,
    "has_next": true
  }
}
```

#### Error Response

```json
{
  "success": false,
  "message": "Invalid email or password",
  "error": null
}
```

### 👨‍💻 Author
#### Savanyv
- Github: @savanyv