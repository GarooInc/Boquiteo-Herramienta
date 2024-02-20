package controllers

import (
	"Boquiteo-Backend/responses"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ReceiveWebhook receives the webhook from Shopify or Resttouch
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

	// TODO implement logging

	var orderNumber int
	if jsonOrderNumber, ok := jsonBody["order_number"]; ok {
		orderNumber = int(jsonOrderNumber.(float64))
	} else {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Order number not found",
			Data:    nil,
		})
		return

	}
	fmt.Println("Order number: ", orderNumber)

	var totalPrice string
	if jsonTotalPrice, ok := jsonBody["total_price"]; ok {
		totalPrice = jsonTotalPrice.(string)
	} else {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Total price not found",
			Data:    nil,
		})
		return
	}
	fmt.Println("Total price: " + totalPrice)

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Webhook received",
		Data:    nil,
	})

	var customer string
	// get jsonBody["customer"]["first_name"] and jsonBody["customer"]["last_name"] and concatenate them to customer
	if jsonFirstName, ok := jsonBody["customer"].(map[string]interface{})["first_name"]; ok {
		customer = jsonFirstName.(string)
	} else {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Customer first name not found",
			Data:    nil,
		})
		return

	}
	if jsonLastName, ok := jsonBody["customer"].(map[string]interface{})["last_name"]; ok {
		customer += " " + jsonLastName.(string)
	} else {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Customer last name not found",
			Data:    nil,
		})
		return
	}
	fmt.Println("Customer: " + customer)
}
