package main

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/api"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	app := fiber.New()

	v := validator.New()
	var j jsoniter.API

	var w sync.WaitGroup

	w.Add(2)

	go iniciarValidate(v, &w)
	go iniciarJson(&j, &w)

	db, err := banco.NewMysqlConn(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	if err != nil {
		log.Fatal(err)
	}

	w.Wait()

	w.Add(1)
	go func() {
		api.BuildRoutes(app, db, v, &j)
		w.Done()
	}()

	var port strings.Builder
	port.WriteString(":")

	if p := os.Getenv("PORT"); p != "" {
		port.WriteString(p)
	} else {
		port.WriteString("8080")
	}

	w.Wait()
	if err = app.Listen(port.String()); err != nil {
		log.Fatal(err)
	}
}

func iniciarJson(j *jsoniter.API, w *sync.WaitGroup) {
	c := jsoniter.Config{
		MarshalFloatWith6Digits: true,
		SortMapKeys:             false,
		UseNumber:               false,
		DisallowUnknownFields:   true,
		OnlyTaggedField:         true,
		CaseSensitive:           true,
	}.Froze()

	*j = c
	w.Done()
}

func iniciarValidate(v *validator.Validate, w *sync.WaitGroup) {
	w.Done()
}
