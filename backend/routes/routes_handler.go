package routes

import (
	"Boquiteo-Backend/controllers"
	"Boquiteo-Backend/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.StandardResponse{
			Status:  http.StatusOK,
			Message: "Service is up and running!",
			Data:    nil,
		})
	})

	router.POST("/webhook", controllers.ReceiveWebhook)

	router.GET("/orders", controllers.GetCurrentOrders)
	router.GET("/orders/:id", controllers.GetOrderById)

	kitchen := router.Group("/kitchen")
	{
		kitchen.PUT("/orders/items", controllers.UpdateItemReady)
		kitchen.PUT("/orders", controllers.SetOrderStatusKitchen)
	}
}
