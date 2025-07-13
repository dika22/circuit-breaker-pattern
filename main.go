package main

import (
	"circuit-breaker-pattern/handler"
	"circuit-breaker-pattern/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiService := service.NewAPIService()
	apiHandler := handler.NewAPIHandler(apiService)

	r.GET("/external", apiHandler.GetExternalData)

	r.Run(":8080")
}
