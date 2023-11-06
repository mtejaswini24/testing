package middleware

import (
	"context"
	"errors"
	"job-portal-api/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Mid struct {
	a *auth.Auth
}

func NewMid(a *auth.Auth) (Mid, error) {
	// 'a' should not be nil because 'nil' indicates that the 'Auth' object does not exist.
	if a == nil {
		// An error is returned when 'a' is 'nil'.
		return Mid{}, errors.New("auth can't be nil")
	}
	//a new 'Mid' instance is returned with 'a' as a field.
	return Mid{a: a}, nil
}

func (m *Mid) Authenticate(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		//get the current request of context
		traceId, ok := ctx.Value(TraceIdKey).(string)
		//extract traceid from request(give the type string bcz context.Value returns an interface{})
		//if traceid not present
		if !ok {
			log.Error().Msg("Traceid not present in the context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error:": http.StatusText(http.StatusInternalServerError)})
			return
		}
		authHeader := c.Request.Header.Get("Authorization")
		//getting authorization header
		parts := strings.Split(authHeader, " ")
		//authHeader contains two strings one is "Bearer" and another "token"
		// Checking the format of the Authorization header
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			log.Error().Err(err).Str("Trace Id", traceId).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims, err := m.a.ValidateToken(parts[1])
		//check the token for validity and returns claims if it's valid
		if err != nil {
			log.Error().Err(err).Str("Trace Id", traceId).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}
		ctx = context.WithValue(ctx, auth.Key, claims)
		req := c.Request.WithContext(ctx)
		c.Request = req
		next(c)
	}

}
