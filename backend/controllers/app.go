package controllers

import (
	"Boquiteo-Backend/configs"
	"Boquiteo-Backend/models"
	"Boquiteo-Backend/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

// GetCurrentOrders
// @Summary Obtener todas las órdenes actuales
// @Description Obtiene todas las órdenes actuales, excluyendo las canceladas y completadas
// @ID get-current-orders
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.MultiOrderResponse
// @Router /orders [get]
func GetCurrentOrders(c *gin.Context) {
	var orders []models.Order

	collection := configs.GetCollection(configs.DB, "orders")

	cursor, err := collection.Find(c, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching orders." + err.Error(),
			Data:    nil,
		})
		return
	}

	if err = cursor.All(c, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching orders." + err.Error(),
			Data:    nil,
		})
		return
	}

	var ordersFiltered []models.Order

	for _, order := range orders {
		if order.Status != "cancelled" && order.Status != "completed" {
			order.ShopifyDetails = nil // No se envía la información de Shopify
			ordersFiltered = append(ordersFiltered, order)
		}
	}

	c.JSON(http.StatusOK, responses.MultiOrderResponse{
		Status:  http.StatusOK,
		Message: "Orders fetched successfully",
		Data:    ordersFiltered,
	})
}

// GetOrderById
// @Summary Obtener una orden por su id
// @Description Obtiene una orden por su id
// @ID get-order-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} responses.OrderResponse
// @Router /orders/{id} [get]
func GetOrderById(c *gin.Context) {
	// Obtener el id del parámetro, y convertirlo a int64
	orderId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid order id." + err.Error(),
			Data:    nil,
		})
		return
	}

	var order models.Order

	collection := configs.GetCollection(configs.DB, "orders")
	err = collection.FindOne(c, bson.M{"order_number": orderId}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching order." + err.Error(),
			Data:    nil,
		})
		return
	}

	var responseOrder = models.Order{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		TotalPrice:  order.TotalPrice,
		Customer:    order.Customer,
		LineItems:   order.LineItems,
		Address:     order.Address,
	}

	c.JSON(http.StatusOK, responses.OrderResponse{
		Status:  http.StatusOK,
		Message: "Order fetched successfully",
		Data:    responseOrder,
	})
}
