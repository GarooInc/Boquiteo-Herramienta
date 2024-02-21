package routes

import (
	"Boquiteo-Backend/controllers"
	"Boquiteo-Backend/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(router *gin.Engine) {
	api := router.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.StandardResponse{
			Status:  http.StatusOK,
			Message: "Service is up and running!",
			Data:    nil,
		})
	})

	api.POST("/webhook", controllers.ReceiveWebhook)
	api.GET("/orders", controllers.GetCurrentOrders)
}