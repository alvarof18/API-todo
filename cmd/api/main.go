package main

import (
	"log"
	"todo-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	//Configurar rutas en Gin
	router := gin.Default()
	routes.Routes(router)
	router.Run("localhost:8080")
	log.Println("Servidor inicializado")
}
