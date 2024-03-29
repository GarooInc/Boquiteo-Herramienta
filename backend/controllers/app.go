package controllers

import (
	"Boquiteo-Backend/configs"
	"Boquiteo-Backend/models"
	"Boquiteo-Backend/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetCurrentOrders
// @Summary Obtener todas las órdenes actuales
// @Description Obtiene todas las órdenes actuales, excluyendo las canceladas y completadas
// @ID get-current-orders
// @Accept  json
// @Produce  json
// @Param debug query int false "Modo Debug: 1 para obtener todas las órdenes, 0 para obtener las órdenes de las últimas 24 horas"
// @Success 200 {object} responses.MultiOrderResponse
// @Router /orders [get]
func GetCurrentOrders(c *gin.Context) {
	var orders []models.Order

	debug := c.Query("debug")

	var cursor *mongo.Cursor
	var bsonFilter bson.M

	if debug == "1" {
		// Modo Debug: regresar todas las ordenes
		bsonFilter = bson.M{}
	} else {
		// Modo Producción: regresar todas las ordenes de las últimas 24 horas
		currentTime := time.Now().UTC()
		oneDayAgo := currentTime.AddDate(0, 0, -1)
		oneDayAgoPrimitive := primitive.NewDateTimeFromTime(oneDayAgo)

		// Filtrar las órdenes de las últimas 24 horas
		bsonFilter = bson.M{"time_order_confirmed": bson.M{"$gte": oneDayAgoPrimitive}}
	}

	collection := configs.GetCollection(configs.DB, "orders")
	opts := options.Find().SetSort(bson.D{{"order_number", -1}})

	cursor, err := collection.Find(c, bsonFilter, opts)

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
		if order.Status != models.Cancelled {
			order.ShopifyDetails = nil // No se envía la información de Shopify
			ordersFiltered = append(ordersFiltered, order)
		}
	}

	fmt.Println(len(ordersFiltered))

	c.JSON(http.StatusOK, responses.MultiOrderResponse{
		Status:  http.StatusOK,
		Message: "Orders fetched successfully",
		Data:    ordersFiltered,
	})
}

// GetPendingOrders
// @Summary Obtener todas las órdenes pendientes
// @Description (Cocina) Obtiene todas las órdenes con status 'Confirmed' o 'Almost' que solo se mostrarán en la cocina
// @ID get-pending-orders
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.MultiOrderResponse
// @Router /kitchen/orders [get]
func GetPendingOrders(c *gin.Context) {
	// regresar todas las ordenes con status Confirmed o Almost
	var orders []models.Order

	collection := configs.GetCollection(configs.DB, "orders")
	opts := options.Find().SetSort(bson.D{{"order_number", -1}})

	cursor, err := collection.Find(c, bson.M{"status": bson.M{"$in": []string{models.Confirmed, models.Almost}}}, opts)
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

	var filteredOrders []models.Order

	for _, order := range orders {
		order.ShopifyDetails = nil // No se envía la información de Shopify
		filteredOrders = append(filteredOrders, order)
	}

	c.JSON(http.StatusOK, responses.MultiOrderResponse{
		Status:  http.StatusOK,
		Message: "Orders fetched successfully",
		Data:    filteredOrders,
	})
}

// GetWaitingOrders
// @Summary Obtener todas las órdenes esperando al repartidor
// @Description (Delivery) Obtiene todas las órdenes con status 'Done' que solo se mostrarán en la sección de repartidor
// @ID get-waiting-orders
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.MultiOrderResponse
// @Router /delivery/orders [get]
func GetWaitingOrders(c *gin.Context) {
	// regresar todas las ordenes con status Done
	var orders []models.Order

	collection := configs.GetCollection(configs.DB, "orders")
	opts := options.Find().SetSort(bson.D{{"order_number", -1}})

	cursor, err := collection.Find(c, bson.M{"status": models.Done}, opts)
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

	var filteredOrders []models.Order

	for _, order := range orders {
		order.ShopifyDetails = nil // No se envía la información de Shopify
		filteredOrders = append(filteredOrders, order)
	}

	c.JSON(http.StatusOK, responses.MultiOrderResponse{
		Status:  http.StatusOK,
		Message: "Orders fetched successfully",
		Data:    filteredOrders,
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
		ID:                 order.ID,
		OrderNumber:        order.OrderNumber,
		Status:             order.Status,
		TotalPrice:         order.TotalPrice,
		Customer:           order.Customer,
		LineItems:          order.LineItems,
		Address:            order.Address,
		Comments:           order.Comments,
		TimeOrderConfirmed: order.TimeOrderConfirmed,
		TimeOrderFulfilled: order.TimeOrderFulfilled,
		TimeOrderPickup:    order.TimeOrderPickup,
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

		if order.Status == models.Confirmed {
			order.Status = models.Almost
		}

	} else {
		order.LineItems[itemIndex].Status = models.ItemPending

		// Si hay al menos un item Listo, la orden pasa a 'Casi listo'. Si no, a 'Confirmado'
		almost := false
		for _, item := range order.LineItems {
			if item.Status == models.ItemReady {
				almost = true
				break
			}
		}

		if almost {
			order.Status = models.Almost
		} else {
			order.Status = models.Confirmed
		}
	}

	// Actualizar la orden en la base de datos
	_, err = collection.UpdateOne(c, bson.M{"order_number": updateItemRequest.OrderId}, bson.M{"$set": bson.M{"line_items": order.LineItems, "status": order.Status}})

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
		Message: "Item updated successfully to '" + order.LineItems[itemIndex].Status + "', order status: '" + order.Status + "'",
		Data:    nil,
	})
}

type UpdateOrderStatusRequest struct {
	OrderId int  `json:"order_id"`
	Status  bool `json:"status"`
}

// SetOrderStatusKitchen
// @Summary Actualizar el estado de una orden en la cocina
// @Description (Cocina) Actualiza el estado de una orden. True para 'Listo', falso para 'Casi listo' o 'Confirmado' dependiendo del estado de los items.
// @ID set-order-status-kitchen
// @Accept  json
// @Produce  json
// @Param updateOrderStatusRequest body UpdateOrderStatusRequest true "Update Order Status Request"
// @Success 200 {object} responses.StandardResponse
// @Router /kitchen/orders [put]
func SetOrderStatusKitchen(c *gin.Context) {
	var updateOrderStatusRequest UpdateOrderStatusRequest
	err := c.BindJSON(&updateOrderStatusRequest)

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
	err = collection.FindOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching order." + err.Error(),
			Data:    nil,
		})
		return
	}

	if order.Status == models.Completed {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Completed order cannot be updated on kitchen.",
			Data:    nil,
		})
		return
	}

	// status = true -> Done
	// status = false -> Almost / Confirmed
	if updateOrderStatusRequest.Status {
		// Verificar que todos los items estén listos
		var itemsNotReady []models.OrderItem
		for _, item := range order.LineItems {
			if item.Status != models.ItemReady {
				itemsNotReady = append(itemsNotReady, item)
			}
		}

		if len(itemsNotReady) > 0 {
			// show in message how many items are not ready
			c.JSON(http.StatusBadRequest, responses.StandardResponse{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("There are %d items not ready.", len(itemsNotReady)),
				Data:    map[string]interface{}{"items": itemsNotReady},
			})
			return
		}

		order.Status = models.Done
		order.TimeOrderFulfilled = primitive.NewDateTimeFromTime(time.Now().UTC())

		// Actualizar la orden en la base de datos, registrando el tiempo en el que se completó
		_, err = collection.UpdateOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}, bson.M{"$set": bson.M{"status": order.Status, "time_order_fulfilled": order.TimeOrderFulfilled}})

		if err != nil {
			log.Println("Error while updating order." + err.Error())
			c.JSON(http.StatusInternalServerError, responses.StandardResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error while updating order." + err.Error(),
				Data:    nil,
			})
			return
		}

	} else {
		// Si hay al menos un item Listo, la orden pasa a 'Casi listo'. Si no, a 'Confirmado'
		almost := false
		for _, item := range order.LineItems {
			if item.Status == models.ItemReady {
				almost = true
				break
			}
		}

		if almost {
			order.Status = models.Almost
		} else {
			order.Status = models.Confirmed
		}

		// Actualizar la orden en la base de datos, eliminando el campo time_order_fulfilled que se vuelve inválido
		_, err = collection.UpdateOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}, bson.M{"$unset": bson.M{"time_order_fulfilled": ""}})

		if err != nil {
			log.Println("Error while updating order." + err.Error())
			c.JSON(http.StatusInternalServerError, responses.StandardResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error while updating order." + err.Error(),
				Data:    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Order updated successfully to '" + order.Status + "'",
		Data:    nil,
	})
}

// SetOrderStatusDelivery
// @Summary Actualizar el estado de una orden en repartidor
// @Description (Repartidor) Actualiza el estado de una orden. True para 'Orden completada', falso para 'Esperando al repartidor'.
// @ID set-order-status-delivery
// @Accept  json
// @Produce  json
// @Param updateOrderStatusRequest body UpdateOrderStatusRequest true "Update Order Status Request"
// @Success 200 {object} responses.StandardResponse
// @Router /delivery/orders [put]
func SetOrderStatusDelivery(c *gin.Context) {
	var updateOrderStatusRequest UpdateOrderStatusRequest
	err := c.BindJSON(&updateOrderStatusRequest)

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
	err = collection.FindOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while fetching order." + err.Error(),
			Data:    nil,
		})
		return
	}

	if order.Status == models.Confirmed || order.Status == models.Almost {
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Order must be 'Done' before updating from Delivery.",
			Data:    nil,
		})
		return

	}

	// status = true -> Delivering
	// status = false -> Done
	if updateOrderStatusRequest.Status {
		order.Status = models.Completed
		order.TimeOrderPickup = primitive.NewDateTimeFromTime(time.Now().UTC())

		// Actualizar la orden en la base de datos
		_, err = collection.UpdateOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}, bson.M{"$set": bson.M{"status": order.Status, "time_order_pickup": order.TimeOrderPickup}})

		if err != nil {
			log.Println("Error while updating order." + err.Error())
			c.JSON(http.StatusInternalServerError, responses.StandardResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error while updating order." + err.Error(),
				Data:    nil,
			})
			return
		}
	} else {
		order.Status = models.Done

		// Actualizar la orden en la base de datos, eliminando el campo time_order_pickup que se vuelve inválido
		_, err = collection.UpdateOne(c, bson.M{"order_number": updateOrderStatusRequest.OrderId}, bson.M{"$unset": bson.M{"time_order_pickup": ""}})

		if err != nil {
			log.Println("Error while updating order." + err.Error())
			c.JSON(http.StatusInternalServerError, responses.StandardResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error while updating order." + err.Error(),
				Data:    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Order updated successfully to '" + order.Status + "'",
		Data:    nil,
	})
}
