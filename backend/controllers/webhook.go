package controllers

import (
	"Boquiteo-Backend/configs"
	"Boquiteo-Backend/models"
	"Boquiteo-Backend/responses"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ReceiveWebhook receives the webhook from Shopify or Resttouch
func ReceiveWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error while reading the request body. " + err.Error())
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Error while reading the request body",
			Data:    nil,
		})
	}

	var jsonBody map[string]interface{}

	if err := json.Unmarshal(body, &jsonBody); err != nil {
		log.Println("Error while parsing the request body. " + err.Error())
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
		log.Println("Order number not found")
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Order number not found",
			Data:    nil,
		})
		return
	}

	var totalPrice float64
	if jsonTotalPrice, ok := jsonBody["total_price"]; ok {
		totalPrice, err = strconv.ParseFloat(jsonTotalPrice.(string), 64)
		if err != nil {
			log.Println("Error while parsing the total price. " + err.Error())
			c.JSON(http.StatusBadRequest, responses.StandardResponse{
				Status:  http.StatusBadRequest,
				Message: "Error while parsing the total price. " + err.Error(),
				Data:    nil,
			})
			return

		}
	} else {
		log.Println("Total price not found")
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Total price not found",
			Data:    nil,
		})
		return
	}

	var customer string
	if jsonFirstName, ok := jsonBody["customer"].(map[string]interface{})["first_name"]; ok {
		customer = jsonFirstName.(string)
	} else {
		log.Println("Customer first name not found")
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
		log.Println("Customer last name not found")
		c.JSON(http.StatusBadRequest, responses.StandardResponse{
			Status:  http.StatusBadRequest,
			Message: "Customer last name not found",
			Data:    nil,
		})
		return
	}

	var lineItems []models.OrderItem // [{"name": "item1", "quantity": 1, "price": 10.00, "vendor": "vendor1"}, ...]
	if jsonLineItems, ok := jsonBody["line_items"]; ok {
		for i, item := range jsonLineItems.([]interface{}) {
			var itemName, itemVendor string
			var itemQuantity, itemPrice float64

			itemNameInterface, itemNameExists := item.(map[string]interface{})["title"]            // String
			itemQuantityInterface, itemQuantityExists := item.(map[string]interface{})["quantity"] // Float64
			itemVendorInterface, itemVendorExists := item.(map[string]interface{})["vendor"]       // String
			itemPriceInterface, itemPriceExists := item.(map[string]interface{})["price"]          // String

			if !itemNameExists || !itemQuantityExists || !itemVendorExists || !itemPriceExists {
				log.Println("Item " + strconv.Itoa(i) + " is missing some fields.")
				c.JSON(http.StatusBadRequest, responses.StandardResponse{
					Status:  http.StatusBadRequest,
					Message: "Item " + strconv.Itoa(i) + " is missing some fields.",
					Data:    nil,
				})
				return
			}

			itemName = itemNameInterface.(string)
			itemQuantity = itemQuantityInterface.(float64)
			itemVendor = itemVendorInterface.(string)
			itemPrice, err = strconv.ParseFloat(itemPriceInterface.(string), 64)

			if err != nil {
				log.Println("Error while parsing the item price. " + err.Error())
				c.JSON(http.StatusBadRequest, responses.StandardResponse{
					Status:  http.StatusBadRequest,
					Message: "Error while parsing the item price. " + err.Error(),
					Data:    nil,
				})
				return
			}

			item := models.OrderItem{
				Item:     i,
				Name:     itemName,
				Quantity: itemQuantity,
				Price:    itemPrice,
				Vendor:   itemVendor,
				Status:   models.ItemPending,
			}

			lineItems = append(lineItems, item)
		}
	}

	var address string
	// get jsonBody["customer"]["default_address"]
	if jsonDefaultAddress, ok := jsonBody["customer"].(map[string]interface{})["default_address"]; ok {
		// if exists, get address1 and address2 and concatenate them to address
		if jsonAddress1, ok := jsonDefaultAddress.(map[string]interface{})["address1"]; ok {
			// strip the address of any commas and blank spaces
			address = jsonAddress1.(string)
		}
		if jsonAddress2, ok := jsonDefaultAddress.(map[string]interface{})["address2"]; ok {
			address += ", " + jsonAddress2.(string)
		}
		if jsonCity, ok := jsonDefaultAddress.(map[string]interface{})["city"]; ok {
			address += ". " + jsonCity.(string)
		}
	}
	address = strings.ReplaceAll(address, "  ", " ") // replace double spaces with single space

	var orderConfirmedTime = primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))

	// Insertar a la base de datos
	var newOrder models.Order

	newOrder.OrderNumber = orderNumber
	newOrder.Status = models.Confirmed
	newOrder.TotalPrice = totalPrice
	newOrder.Customer = customer
	newOrder.LineItems = lineItems
	newOrder.Address = address
	newOrder.TimeOrderConfirmed = orderConfirmedTime
	newOrder.ShopifyDetails = jsonBody

	// Insertar a la base de datos
	collection := configs.GetCollection(configs.DB, "orders")
	insertResult, err := collection.InsertOne(c, newOrder)
	if err != nil {
		log.Println("Error while inserting the order to the database." + err.Error())
		c.JSON(http.StatusInternalServerError, responses.StandardResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error while inserting the order to the database." + err.Error(),
			Data:    nil,
		})
		return
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Order inserted to the database",
		Data: map[string]interface{}{
			"order_id": insertResult.InsertedID,
		},
	})
}
