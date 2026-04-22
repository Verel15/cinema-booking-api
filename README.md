# Cinema Booking API

## Table of Contents
- [Project Structure](#project-structure)
- [Authentication](#authentication)
  - [Register](#register)
  - [Login](#login)
  - [Google OAuth](#google-oauth)
  - [Refresh Token](#refresh-token)
- [Authorization (RBAC)](#authorization-rbac)
  - [Roles](#roles)
  - [Guard Middleware](#guard-middleware)
  - [Permission Configuration](#permission-configuration)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)

---

## Project Structure

```
cinema-booking-api/
├── cmd/app/
│   └── main.go                 # Application entry point
├── internal/
│   ├── auth/                   # Authentication module
│   │   ├── delivery/http/      # HTTP handlers
│   │   ├── dto/                # Data transfer objects
│   │   ├── repository/         # Database repository
│   │   └── usecase/            # Business logic
│   ├── user/                   # User module
│   │   ├── domain/             # Domain models & interfaces
│   │   ├── delivery/http/      # HTTP handlers
│   │   ├── dto/                # Data transfer objects
│   │   ├── repository/         # Database repository
│   │   └── usecase/            # Business logic
│   ├── movie/                  # Movie module
│   └── database/               # Database connection
├── pkg/
│   ├── enums/                  # Enumerations
│   ├── jwt/                    # JWT token handling
│   ├── middleware/             # Middlewares (Auth, RBAC)
│   ├── pagination/             # Pagination utility
│   ├── response/               # Response formatting
│   └── utils/                  # Utility functions
└── configs/
```

---

## Authentication

### How It Works

The authentication system uses **JWT (JSON Web Tokens)** with two types of tokens:

1. **Access Token** - Short-lived (1 hour), used for API requests
2. **Refresh Token** - Long-lived (7 days), used to get new access tokens

### Register

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user",
    "status": "active",
    "provider": "email"
  }
}
```

### Login

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

### Google OAuth

**Step 1: Get Google Auth URL**

**Endpoint:** `GET /api/v1/auth/google`

**Response:**
```json
{
  "url": "https://accounts.google.com/o/oauth2/v2/auth?..."
}
```

Redirect user to this URL to login with Google.

**Step 2: Google Callback**

**Endpoint:** `GET /api/v1/auth/google/callback?code=...`

Google will redirect back with an authorization code. Exchange it for user info and create/login the user.

### Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** New access token and refresh token pair.

---

## Authorization (RBAC)

### Roles

| Role | Description |
|------|-------------|
| `admin` | Full access to all resources |
| `user` | Limited access based on permissions |

### Guard Middleware

The Guard middleware checks if the authenticated user has the required role to access a specific endpoint.

**How it works:**

1. **Authentication** - First, the user must be authenticated via JWT token
2. **Role Extraction** - The user's role is extracted from the token
3. **Permission Check** - The Guard checks if the user's role is allowed for the requested path
4. **Access Control** - If allowed, proceed; otherwise, return 403 Forbidden

**Code Flow:**

```
Request → AuthMiddleware (validate token) → RBAC.Guard() (check role) → Handler
```

### Permission Configuration

Permissions are registered in `main.go`:

```go
// Initialize RBAC
rbac := middleware.NewRBAC()

// Register permissions
rbac.RegisterPermission("/api/v1/movies", middleware.RoleAdmin)
rbac.RegisterPermission("/api/v1/movies/:id", middleware.RoleAdmin)
rbac.RegisterPermission("/api/v1/users", middleware.RoleAdmin)
rbac.RegisterPermission("/api/v1/users/:username", middleware.RoleAdmin, middleware.RoleUser)
```

**Usage in Routes:**

```go
// Admin only routes
adminRoutes := protectedRoutes.Group("/movies")
adminRoutes.Use(rbac.Guard())  // Only admin can access
{
    adminRoutes.POST("/", movieHD.CreateMovie)
    adminRoutes.PUT("/:id", movieHD.UpdateMovie)
    adminRoutes.DELETE("/:id", movieHD.DeleteMovie)
}

// Authenticated users (admin or user)
protectedRoutes.GET("/movies", movieHD.GetAllMovies)
protectedRoutes.GET("/movies/:id", movieHD.GetMovieByID)
```

---

## API Endpoints

### Public Routes (No Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login with email/password |
| GET | `/api/v1/auth/google` | Get Google OAuth URL |
| GET | `/api/v1/auth/google/callback` | Google OAuth callback |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/health` | Health check |

### Protected Routes (Auth Required)

| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/api/v1/movies` | User+ | Get all movies |
| GET | `/api/v1/movies/:id` | User+ | Get movie by ID |
| POST | `/api/v1/movies` | Admin | Create movie |
| PUT | `/api/v1/movies/:id` | Admin | Update movie |
| DELETE | `/api/v1/movies/:id` | Admin | Delete movie |
| GET | `/api/v1/users` | Admin | Get all users |
| GET | `/api/v1/users/:username` | User+ | Get user by username |
| GET | `/api/v1/me` | User+ | Get current user info |

**Access Levels:**
- `Admin` - Admin only
- `User+` - Both Admin and User roles

### Using Protected Routes

Include the access token in the `Authorization` header:

```bash
curl -X GET http://localhost:5050/api/v1/movies \
  -H "Authorization: Bearer <access_token>"
```

---

## Environment Variables

Create a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=cinema_booking

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Google OAuth (optional for development)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:5050/api/v1/auth/google/callback
```

---

## Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/app/main.go
```

The server will start on `http://localhost:5050`