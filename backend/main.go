package main

import (
	"Boquiteo-Backend/configs"
	docs "Boquiteo-Backend/docs"
	"Boquiteo-Backend/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CORS() gin.HandlerFunc {
	// Reference: https://github.com/gin-contrib/cors/issues/29
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORS())

	// MongoDB
	configs.LoadSetup()

	// Swagger
	docs.SwaggerInfo.Title = "Boquiteo API"
	docs.SwaggerInfo.Description = "API para la automatización de órdenes de Boquiteo"
	docs.SwaggerInfo.Version = "0.1.0"
	docs.SwaggerInfo.Host = "boquiteo-garoo.koyeb.app/api"
	docs.SwaggerInfo.BasePath = "/"

	// Routes
	//urlSwagger := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.Routes(r)

	err := r.Run() // listen and serve on
	if err != nil {
		panic(err)
	}
}
