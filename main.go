package main

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/api"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/database"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strings"
)

func main() {
	app := fiber.New()

	db, err := database.NewMysqlConn(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.BuildRoutes(app, &db)

	var port strings.Builder
	port.WriteString(":")

	if p := os.Getenv("PORT"); p != "" {
		port.WriteString(p)
	} else {
		port.WriteString("8080")
	}

	if err = app.Listen(port.String()); err != nil {
		log.Fatal(err)
	}
}
