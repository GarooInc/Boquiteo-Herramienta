package controllers

import (
	"Boquiteo-Backend/responses"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ReceiveWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Error while reading the request body",
			Data:    nil,
		})
	}

	var jsonBody map[string]interface{}

	if err := json.Unmarshal(body, &jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Error while parsing the request body. " + err.Error(),
			Data:    nil,
		})
		return
	}

	if orderNumber, ok := jsonBody["order_number"]; ok {
		fmt.Printf("Order number: %v\n", orderNumber)
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Webhook received",
		Data:    nil,
	})
}
