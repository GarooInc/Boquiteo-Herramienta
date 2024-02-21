package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CONFIRMED  = "Confirmado"              // La orden ingresa al sistema
	ALMOST     = "Casi listo"              // El primer plato de la orden está listo
	DONE       = "Esperando al repartidor" // La orden está lista para ser recogida
	DELIVERING = "En camino"               // La orden fue recogida por el repartidor
	COMPLETED  = "Orden completada"        // La orden fue entregada
	CANCELLED  = "Cancelado"               // La orden fue cancelada
)

// Order struct para el modelo de datos de la orden
type Order struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderNumber        int                `json:"order_number" bson:"order_number"`
	Status             string             `json:"status" bson:"status"`
	TotalPrice         float64            `json:"total_price" bson:"total_price"`
	Customer           string             `json:"customer" bson:"customer"`
	LineItems          []OrderItem        `json:"line_items" bson:"line_items"`
	Address            string             `json:"address" bson:"address"`
	TimeOrderConfirmed primitive.DateTime `json:"time_order_confirmed" bson:"time_order_confirmed"`
	TimeOrderFulfilled primitive.DateTime `json:"time_order_fulfilled,omitempty" bson:"time_order_fulfilled,omitempty"`
	TimeOrderPickup    primitive.DateTime `json:"time_order_pickup,omitempty" bson:"time_order_pickup,omitempty"`
	ShopifyDetails     interface{}        `json:"shopify_details" bson:"shopify_details"`
}

type OrderItem struct {
	Name     string  `json:"name" bson:"name"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
	Vendor   string  `json:"vendor" bson:"vendor"`
}
