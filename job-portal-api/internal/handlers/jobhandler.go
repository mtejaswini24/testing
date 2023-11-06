package handlers

import (
	"encoding/json"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateJob(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var jobData models.NewJob //store the data in dummmy struct
	err := json.NewDecoder(c.Request.Body).Decode(&jobData)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	validate := validator.New()
	err = validate.Struct(jobData)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide job role and Description"})
		return
	}
	stringCmpnyId := c.Param("id")
	cid, err := strconv.ParseUint(stringCmpnyId, 10, 64)
	if err != nil {
		log.Print("conversion string to int error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error found at conversion.."})
		return
	}
	job, err := h.service.CreateJob(ctx, jobData, cid)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in adding job details")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "adding job failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, job)
}
func (h *Handler) FetchJob(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	jobDetails, err := h.service.FetchJob()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in fetching job")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "fetching failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, jobDetails)
}
func (h *Handler) FetchJobById(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	stringJobId := c.Param("id")
	jid, err := strconv.ParseUint(stringJobId, 10, 64)
	if err != nil {
		log.Print("conversion string to int error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error found at conversion.."})
		return
	}
	jobData, err := h.service.FetchJobById(jid)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("job is not there with that id")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "job fetching failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, jobData)
}
func (h *Handler) FetchJobByCompanyId(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	stringCompanyId := c.Param("id")
	cid, err := strconv.ParseUint(stringCompanyId, 10, 64)
	if err != nil {
		log.Print("conversion string to int error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error found at conversion.."})
		return
	}
	job, err := h.service.FetchJobByCompanyId(cid)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("job is not there with that company id")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "job fetching with company failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, job)
}
