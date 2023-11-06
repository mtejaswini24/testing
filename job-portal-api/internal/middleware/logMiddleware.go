package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type key string

const TraceIdKey key = "1"

func (a *Mid) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := uuid.NewString()
		//generate new unique identifier
		ctx := c.Request.Context()
		//get current context from gin
		ctx = context.WithValue(ctx, TraceIdKey, traceId)
		//adding traceid and key for upcoming use
		req := c.Request.WithContext(ctx)
		//create a new copy and updating with ctx contex that contain traceid and key
		c.Request = req
		//replacing the original request with copy req
		log.Info().Str("Traceid", traceId).Str("Method", c.Request.Method).Str("url path", c.Request.URL.Path).Msg("Request started")
		defer log.Info().Str("Trace Id", traceId).Str("Method", c.Request.Method).Str("URL Path", c.Request.URL.Path).Int("status Code", c.Writer.Status()).Msg("Request processing completed")
		c.Next()
	}
}
