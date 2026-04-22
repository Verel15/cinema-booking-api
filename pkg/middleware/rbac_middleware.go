package middleware

import (
	"cinema-booking-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role-based access control
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// Permission defines what roles can access a resource
type Permission struct {
	Roles []Role
}

// RBAC - Role Based Access Control
type RBAC struct {
	permissions map[string]Permission
}

func NewRBAC() *RBAC {
	return &RBAC{
		permissions: make(map[string]Permission),
	}
}

// RegisterPermission registers a permission for a route
func (r *RBAC) RegisterPermission(path string, roles ...Role) {
	r.permissions[path] = Permission{
		Roles: roles,
	}
}

// Middleware to check if user has permission
func (r *RBAC) Guard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by auth middleware)
		user, exists := c.Get("user")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		// Get user role
		userMap, ok := user.(map[string]interface{})
		if !ok {
			response.Error(c, http.StatusForbidden, "invalid user data")
			c.Abort()
			return
		}

		userRole := userMap["role"].(string)

		// Get the full path without query params
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Check if route requires specific permissions
		perm, exists := r.permissions[path]
		if !exists {
			// No specific permission required, allow
			c.Next()
			return
		}

		// Check if user role is allowed
		allowed := false
		for _, role := range perm.Roles {
			if string(role) == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "you don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
