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

// UpdateItemRequest
// @Description Estructura de la solicitud para actualizar el estado de un item
type UpdateItemRequest struct {
	OrderId   int  `json:"order_id"`
	ItemId    int  `json:"item_id"`
	ItemReady bool `json:"item_ready"`
}

// UpdateItemReady
// @Summary Actualizar el estado de un item
// @Description (Cocina) Actualiza el estado de un item, si está listo o no. False para 'Pendiente', true para 'Listo'.
// @ID update-item-ready
// @Accept  json
// @Produce  json
// @Param updateItemRequest body UpdateItemRequest true "Update Item Request"
// @Success 200 {object} responses.StandardResponse
// @Router /kitchen/orders/items [put]
func UpdateItemReady(c *gin.Context) {
	var updateItemRequest UpdateItemRequest
	err := c.BindJSON(&updateItemRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request." + err.Error(),
			Data:    nil,
		})
		return
	}

	// Obtener la orden
	var order models.Order

	collection := configs.GetCollection(configs.DB, "orders")
	err = collection.FindOne(c, bson.M{"order_number": updateItemRequest.OrderId}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching order." + err.Error(),
			Data:    nil,
		})
		return
	}

	// Encontrar el item
	var itemIndex = -1
	for i, item := range order.LineItems {
		if item.Item == updateItemRequest.ItemId {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		c.JSON(http.StatusNotFound, responses.StandardResponse{
			Status:  http.StatusNotFound,
			Message: "Item not found in the order.",
			Data:    nil,
		})
		return
	}

	// Actualizar el item
	if updateItemRequest.ItemReady {
		order.LineItems[itemIndex].Status = models.ItemReady
	} else {
		order.LineItems[itemIndex].Status = models.ItemPending
	}

	// Actualizar la orden en la base de datos
	_, err = collection.UpdateOne(c, bson.M{"order_number": updateItemRequest.OrderId}, bson.M{"$set": bson.M{"line_items": order.LineItems}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while updating order." + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Item updated successfully to '" + order.LineItems[itemIndex].Status + "'",
		Data:    nil,
	})
}
