package main

import (
	"log"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	//Cargar configuraciones
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//Configurar Database
	db, errDB := database.Connected(*cfg)
	if errDB != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	//Configurar rutas en Gin
	router := gin.Default()
	routes.Routes(router, db)
	log.Println("Servidor inicializado")
	router.Run(cfg.Host + ":" + cfg.HostPort)

}
