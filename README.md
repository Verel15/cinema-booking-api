# Cinema Booking API

## Table of Contents
- [Project Structure](#project-structure)
- [How to Add a New API Endpoint](#how-to-add-a-new-api-endpoint)
  - [Overview](#overview)
  - [Step 1 — Domain (Entity + Interfaces)](#step-1--domain-entity--interfaces)
  - [Step 2 — DTO (Request / Response)](#step-2--dto-request--response)
  - [Step 3 — Repository (Database Layer)](#step-3--repository-database-layer)
  - [Step 4 — Usecase (Business Logic)](#step-4--usecase-business-logic)
  - [Step 5 — Handler (HTTP Layer)](#step-5--handler-http-layer)
  - [Step 6 — Wire Everything in main.go](#step-6--wire-everything-in-maingo)
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

---

## How to Add a New API Endpoint

### Overview

โปรเจกต์นี้ใช้สถาปัตยกรรมแบบ **Clean Architecture** แบ่งออกเป็น 4 Layer หลัก ทุก module (เช่น `movie`, `user`) มีโครงสร้างเดียวกัน

```
Request
  │
  ▼
Handler          ← รับ HTTP Request, แปลง JSON, ส่งต่อให้ Usecase
  │
  ▼
Usecase          ← Business Logic, orchestrate การทำงาน
  │
  ▼
Repository       ← คุยกับ Database ผ่าน GORM
  │
  ▼
Database (PostgreSQL)
```

**ลำดับการเขียนไฟล์ที่ถูกต้อง:**

```
1. domain/       → นิยาม Entity + Interface (Repository & Usecase)
2. dto/          → นิยาม Request / Response struct
3. repository/   → implement Repository interface (GORM)
4. usecase/      → implement Usecase interface (business logic)
5. delivery/http/→ implement HTTP Handler
6. main.go       → wire ทุกอย่างเข้าด้วยกัน + ลงทะเบียน route
```

ตัวอย่างด้านล่างใช้ module `screening` (รอบฉาย) เป็น reference

---

### Step 1 — Domain (Entity + Interfaces)

สร้างโฟลเดอร์ `internal/screening/domain/`

**`screening.go`** — นิยาม struct ของ Entity ที่จะ map กับตาราง DB

```go
package domain

import "time"

type Screening struct {
    ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    MovieID   string    `gorm:"not null" json:"movie_id"`
    StartsAt  time.Time `gorm:"not null" json:"starts_at"`
    Hall      string    `gorm:"not null" json:"hall"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**`repository.go`** — นิยาม Interface ของ Repository (บอกว่า DB ต้องทำอะไรได้บ้าง)

```go
package domain

type ScreeningRepository interface {
    Create(s *Screening) error
    GetAll() ([]Screening, error)
    GetByID(id string) (*Screening, error)
    Delete(id string) error
}
```

**`usecase.go`** — นิยาม Interface ของ Usecase (บอกว่า Business Logic ต้องทำอะไรได้บ้าง)

```go
package domain

import "cinema-booking-api/internal/screening/dto"

type ScreeningUsecase interface {
    CreateScreening(req dto.CreateScreeningRequest) (*dto.ScreeningResponse, error)
    GetAllScreenings() ([]dto.ScreeningResponse, error)
    GetScreeningByID(id string) (*dto.ScreeningResponse, error)
    DeleteScreening(id string) error
}
```

> **หลักการ:** Domain layer ไม่รู้จัก layer อื่น ทำให้ swap implementation ได้ง่าย (เช่น เปลี่ยนจาก GORM เป็น raw SQL โดยไม่กระทบ Usecase หรือ Handler)

---

### Step 2 — DTO (Request / Response)

สร้างไฟล์ `internal/screening/dto/screening_dto.go`

```go
package dto

import "time"

// ใช้รับข้อมูลจาก client
type CreateScreeningRequest struct {
    MovieID  string    `json:"movie_id"  binding:"required,uuid"`
    StartsAt time.Time `json:"starts_at" binding:"required"`
    Hall     string    `json:"hall"      binding:"required"`
}

// ใช้ส่งข้อมูลกลับให้ client
type ScreeningResponse struct {
    ID       string    `json:"id"`
    MovieID  string    `json:"movie_id"`
    StartsAt time.Time `json:"starts_at"`
    Hall     string    `json:"hall"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

> **หลักการ:** DTO แยก struct ที่รับ/ส่งกับ client ออกจาก Entity ของ DB เพื่อความปลอดภัยและ flexibility

---

### Step 3 — Repository (Database Layer)

สร้างไฟล์ `internal/screening/repository/gorm_screening_repository.go`

```go
package repository

import (
    "cinema-booking-api/internal/screening/domain"
    "gorm.io/gorm"
)

type gormScreeningRepository struct {
    db *gorm.DB
}

// NewScreeningRepository คืน interface ไม่ใช่ concrete type
func NewScreeningRepository(db *gorm.DB) domain.ScreeningRepository {
    return &gormScreeningRepository{db: db}
}

func (r *gormScreeningRepository) Create(s *domain.Screening) error {
    return r.db.Create(s).Error
}

func (r *gormScreeningRepository) GetAll() ([]domain.Screening, error) {
    var screenings []domain.Screening
    err := r.db.Find(&screenings).Error
    return screenings, err
}

func (r *gormScreeningRepository) GetByID(id string) (*domain.Screening, error) {
    var s domain.Screening
    err := r.db.First(&s, "id = ?", id).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &s, nil
}

func (r *gormScreeningRepository) Delete(id string) error {
    return r.db.Delete(&domain.Screening{}, "id = ?", id).Error
}
```

> **หลักการ:** Repository layer รับผิดชอบแค่การ query/write DB เท่านั้น ไม่มี business logic

---

### Step 4 — Usecase (Business Logic)

สร้างไฟล์ `internal/screening/usecase/screening_usecase.go`

```go
package usecase

import (
    "cinema-booking-api/internal/screening/domain"
    "cinema-booking-api/internal/screening/dto"
)

type screeningUsecase struct {
    repo domain.ScreeningRepository
}

// NewScreeningUsecase รับ interface ไม่ใช่ concrete type → testable
func NewScreeningUsecase(repo domain.ScreeningRepository) domain.ScreeningUsecase {
    return &screeningUsecase{repo: repo}
}

func (u *screeningUsecase) CreateScreening(req dto.CreateScreeningRequest) (*dto.ScreeningResponse, error) {
    s := &domain.Screening{
        MovieID:  req.MovieID,
        StartsAt: req.StartsAt,
        Hall:     req.Hall,
    }
    if err := u.repo.Create(s); err != nil {
        return nil, err
    }
    return u.mapToResponse(s), nil
}

func (u *screeningUsecase) GetAllScreenings() ([]dto.ScreeningResponse, error) {
    screenings, err := u.repo.GetAll()
    if err != nil {
        return nil, err
    }
    responses := make([]dto.ScreeningResponse, len(screenings))
    for i, s := range screenings {
        responses[i] = *u.mapToResponse(&s)
    }
    return responses, nil
}

func (u *screeningUsecase) GetScreeningByID(id string) (*dto.ScreeningResponse, error) {
    s, err := u.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if s == nil {
        return nil, nil
    }
    return u.mapToResponse(s), nil
}

func (u *screeningUsecase) DeleteScreening(id string) error {
    return u.repo.Delete(id)
}

// mapToResponse แปลง Entity → DTO ที่ส่งกลับ client
func (u *screeningUsecase) mapToResponse(s *domain.Screening) *dto.ScreeningResponse {
    return &dto.ScreeningResponse{
        ID:        s.ID,
        MovieID:   s.MovieID,
        StartsAt:  s.StartsAt,
        Hall:      s.Hall,
        CreatedAt: s.CreatedAt,
        UpdatedAt: s.UpdatedAt,
    }
}
```

---

### Step 5 — Handler (HTTP Layer)

สร้างไฟล์ `internal/screening/delivery/http/screening_handler.go`

```go
package http

import (
    "cinema-booking-api/internal/screening/domain"
    "cinema-booking-api/internal/screening/dto"
    "cinema-booking-api/pkg/response"
    "net/http"

    "github.com/gin-gonic/gin"
)

type ScreeningHandler struct {
    usecase domain.ScreeningUsecase
}

func NewScreeningHandler(u domain.ScreeningUsecase) *ScreeningHandler {
    return &ScreeningHandler{usecase: u}
}

func (h *ScreeningHandler) CreateScreening(c *gin.Context) {
    var req dto.CreateScreeningRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    res, err := h.usecase.CreateScreening(req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, http.StatusCreated, res)
}

func (h *ScreeningHandler) GetAllScreenings(c *gin.Context) {
    screenings, err := h.usecase.GetAllScreenings()
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(c, http.StatusOK, screenings)
}

func (h *ScreeningHandler) GetScreeningByID(c *gin.Context) {
    id := c.Param("id")

    res, err := h.usecase.GetScreeningByID(id)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    if res == nil {
        response.Error(c, http.StatusNotFound, "screening not found")
        return
    }

    response.Success(c, http.StatusOK, res)
}

func (h *ScreeningHandler) DeleteScreening(c *gin.Context) {
    id := c.Param("id")

    if err := h.usecase.DeleteScreening(id); err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, http.StatusOK, gin.H{"message": "screening deleted successfully"})
}
```

> **หลักการ:** Handler รับผิดชอบแค่ HTTP เท่านั้น — parse request, เรียก usecase, format response ไม่มี business logic ใน handler

---

### Step 6 — Wire Everything in main.go

เปิดไฟล์ [cmd/app/main.go](cmd/app/main.go) แล้วเพิ่มในสามจุด

**1. Import**
```go
import (
    // ... existing imports ...
    screeningDomain "cinema-booking-api/internal/screening/domain"
    screeningHandler "cinema-booking-api/internal/screening/delivery/http"
    screeningRepo "cinema-booking-api/internal/screening/repository"
    screeningUsecase "cinema-booking-api/internal/screening/usecase"
)
```

**2. Auto-migrate Entity**
```go
if err := database.Migrate(db,
    &movieDomain.Movie{},
    &userDomain.User{},
    &screeningDomain.Screening{}, // ← เพิ่มตรงนี้
); err != nil {
    log.Fatalf("Database migration failed: %v", err)
}
```

**3. Wire + Register Routes**
```go
// Initialize Screening Module
screeningRepository := screeningRepo.NewScreeningRepository(db)
screeningUC := screeningUsecase.NewScreeningUsecase(screeningRepository)
screeningHD := screeningHandler.NewScreeningHandler(screeningUC)

// Register Routes
screeningRoutes := api.Group("/screenings")
{
    screeningRoutes.GET("/", screeningHD.GetAllScreenings)
    screeningRoutes.GET("/:id", screeningHD.GetScreeningByID)
}

// Protected routes (ถ้าต้องการ auth)
protectedScreeningRoutes := protectedRoutes.Group("/screenings")
{
    protectedScreeningRoutes.POST("/", screeningHD.CreateScreening)
    protectedScreeningRoutes.DELETE("/:id", screeningHD.DeleteScreening)
}
```

---

### สรุปไฟล์ที่ต้องสร้างทั้งหมด

```
internal/screening/
├── domain/
│   ├── screening.go          # Entity struct
│   ├── repository.go         # Repository interface
│   └── usecase.go            # Usecase interface
├── dto/
│   └── screening_dto.go      # Request / Response structs
├── repository/
│   └── gorm_screening_repository.go   # GORM implementation
├── usecase/
│   └── screening_usecase.go  # Business logic
└── delivery/http/
    └── screening_handler.go  # HTTP Handler
```

แก้ไข 1 ไฟล์:
```
cmd/app/main.go               # Wire + Register route
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