package main

import (
	db "demo-go-firebase/firebase"
	"demo-go-firebase/router"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Failed to load .env file")
	}

	db.InitFirestoreCommunicator()
	defer db.CloseFirestoreCommunicator()

	app := fiber.New()
	router.InitializeRoutes(app)
	runService(app)
}

func runService(app *fiber.App) {
	fmt.Println("Init app successfully connectat port : " + getPort())
	app.Listen("0.0.0.0" + getPort())
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
