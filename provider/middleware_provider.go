package provider

import (
	"dinacom-11.0-backend/middleware"

	"github.com/gin-gonic/gin"
)

type MiddlewareProvider interface {
	ProvideAuthMiddleware() gin.HandlerFunc
}

type middlewareProvider struct {
	authMiddleware gin.HandlerFunc
}

func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {
	return &middlewareProvider{
		authMiddleware: middleware.AuthMiddleware(),
	}
}

func (m *middlewareProvider) ProvideAuthMiddleware() gin.HandlerFunc {
	return m.authMiddleware
}
