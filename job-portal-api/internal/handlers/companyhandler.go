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

func (h *Handler) CreateCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var companyData models.NewCompany //store the data in dummmy struct
	err := json.NewDecoder(c.Request.Body).Decode(&companyData)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	validate := validator.New()
	err = validate.Struct(companyData)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide companyname and location"})
		return
	}
	companyDetails, err := h.service.CreateCompany(ctx, companyData)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in adding company")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "company adding failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, companyDetails)
}

func (h *Handler) FetchCompanies(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	companyDetails, err := h.service.FetchCompanies()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in fetching company details")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "getting companydetails failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, companyDetails)
}
func (h *Handler) FetchCompanyById(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	stringCmpnyId := c.Param("id")
	cid, err := strconv.ParseUint(stringCmpnyId, 10, 64)
	if err != nil {

		log.Print("conversion string to int error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error found at conversion.."})
		return

	}
	companyData, err := h.service.FetchCompanyById(cid)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in fetching company by id")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "failed to get company details by id"})
		return
	}
	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, companyData)
}
