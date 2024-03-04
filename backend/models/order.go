package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	Confirmed   = "Confirmado"              // La orden ingresa al sistema
	Almost      = "Casi listo"              // El primer plato de la orden está listo
	Done        = "Esperando al repartidor" // La orden está lista para ser recogida
	Delivering  = "En camino"               // La orden fue recogida por el repartidor
	Completed   = "Orden completada"        // La orden fue entregada
	Cancelled   = "Cancelado"               // La orden fue cancelada
	ItemPending = "Pendiente"               // El item de la orden está pendiente
	ItemReady   = "Listo"                   // El item de la orden está listo
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
	Comments           string             `json:"comments,omitempty" bson:"comments,omitempty"`
	TimeOrderConfirmed primitive.DateTime `json:"time_order_confirmed" bson:"time_order_confirmed" swaggertype:"string" format:"date-time"`
	TimeOrderFulfilled primitive.DateTime `json:"time_order_fulfilled,omitempty" bson:"time_order_fulfilled,omitempty" swaggertype:"string" format:"date-time"`
	TimeOrderPickup    primitive.DateTime `json:"time_order_pickup,omitempty" bson:"time_order_pickup,omitempty" swaggertype:"string" format:"date-time"`
	ShopifyDetails     interface{}        `json:"shopify_details,omitempty" bson:"shopify_details"`
}

type OrderItem struct {
	Item     int      `json:"item" bson:"item"` // ID del item dentro de la orden
	Name     string   `json:"name" bson:"name"`
	Quantity float64  `json:"quantity" bson:"quantity"`
	Price    float64  `json:"price" bson:"price"`
	Vendor   string   `json:"vendor" bson:"vendor"`
	Status   string   `json:"status" bson:"status"`
	Variant  string   `json:"variant,omitempty" bson:"variant,omitempty"`
	Options  []string `json:"options,omitempty" bson:"options,omitempty"`
}
