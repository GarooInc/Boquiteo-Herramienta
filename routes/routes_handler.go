package routes

import (
	"Boquiteo-Backend/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.StandardResponse{
			Status:  http.StatusOK,
			Message: "Service is up and running!",
			Data:    nil,
		})
	})
}
