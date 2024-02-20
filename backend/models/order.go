package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Order struct
type Order struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderNumber        int                `json:"order_number" bson:"order_number"`
	TotalPrice         string             `json:"total_price" bson:"total_price"`
	Customer           string             `json:"customer" bson:"customer"`
	LineItems          []interface{}      `json:"line_items" bson:"line_items"`
	Address            string             `json:"address" bson:"address"`
	TimeOrderConfirmed primitive.DateTime `json:"time_order_confirmed" bson:"time_order_confirmed"`
	ShopifyDetails     interface{}        `json:"shopify_details" bson:"shopify_details"`
}
