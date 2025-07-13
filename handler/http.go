package handler

import (
	"net/http"

	"circuit-breaker-pattern/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	Service service.APIService
}

func NewAPIHandler(service service.APIService) *APIHandler {
	return &APIHandler{
		Service: service,
	}
}

func (h *APIHandler) GetExternalData(c *gin.Context) {
	data, err := h.Service.GetExternalData()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Failed to fetch from external API",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
