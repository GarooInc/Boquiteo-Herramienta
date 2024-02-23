// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/kitchen/orders": {
            "get": {
                "description": "(Cocina) Obtiene todas las órdenes con status 'Confirmed' o 'Almost' que solo se mostrarán en la cocina",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Obtener todas las órdenes pendientes",
                "operationId": "get-pending-orders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.MultiOrderResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "(Cocina) Actualiza el estado de una orden. True para 'Listo', falso para 'Casi listo' o 'Confirmado' dependiendo del estado de los items.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Actualizar el estado de una orden en la cocina",
                "operationId": "set-order-status-kitchen",
                "parameters": [
                    {
                        "description": "Update Order Status Request",
                        "name": "updateOrderStatusRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateOrderStatusRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.StandardResponse"
                        }
                    }
                }
            }
        },
        "/kitchen/orders/items": {
            "put": {
                "description": "(Cocina) Actualiza el estado de un item, si está listo o no. False para 'Pendiente', true para 'Listo'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Actualizar el estado de un item",
                "operationId": "update-item-ready",
                "parameters": [
                    {
                        "description": "Update Item Request",
                        "name": "updateItemRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.StandardResponse"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "description": "Obtiene todas las órdenes actuales, excluyendo las canceladas y completadas",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Obtener todas las órdenes actuales",
                "operationId": "get-current-orders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.MultiOrderResponse"
                        }
                    }
                }
            }
        },
        "/orders/{id}": {
            "get": {
                "description": "Obtiene una orden por su id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Obtener una orden por su id",
                "operationId": "get-order-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.OrderResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.UpdateItemRequest": {
            "description": "Estructura de la solicitud para actualizar el estado de un item",
            "type": "object",
            "properties": {
                "item_id": {
                    "type": "integer"
                },
                "item_ready": {
                    "type": "boolean"
                },
                "order_id": {
                    "type": "integer"
                }
            }
        },
        "controllers.UpdateOrderStatusRequest": {
            "type": "object",
            "properties": {
                "order_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "models.Order": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "address": {
                    "type": "string"
                },
                "customer": {
                    "type": "string"
                },
                "line_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.OrderItem"
                    }
                },
                "order_number": {
                    "type": "integer"
                },
                "shopify_details": {},
                "status": {
                    "type": "string"
                },
                "time_order_confirmed": {
                    "type": "string",
                    "format": "date-time"
                },
                "time_order_fulfilled": {
                    "type": "string",
                    "format": "date-time"
                },
                "time_order_pickup": {
                    "type": "string",
                    "format": "date-time"
                },
                "total_price": {
                    "type": "number"
                }
            }
        },
        "models.OrderItem": {
            "type": "object",
            "properties": {
                "item": {
                    "description": "ID del item dentro de la orden",
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "number"
                },
                "status": {
                    "type": "string"
                },
                "vendor": {
                    "type": "string"
                }
            }
        },
        "responses.MultiOrderResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Order"
                    }
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "responses.OrderResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.Order"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "responses.StandardResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
