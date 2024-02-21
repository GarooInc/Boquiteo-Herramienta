package responses

import "Boquiteo-Backend/models"

type StandardResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type OrderResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []models.Order `json:"data"`
}
