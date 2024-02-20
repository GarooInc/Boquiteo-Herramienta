package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Order struct
type Order struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderNumber    int                `json:"order_number" bson:"order_number"`
	TotalPrice     string             `json:"total_price" bson:"total_price"`
	Customer       string             `json:"customer" bson:"customer"`
	ShopifyDetails interface{}        `json:"shopify_details" bson:"shopify_details"`
}
