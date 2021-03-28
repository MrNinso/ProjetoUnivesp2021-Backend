package main

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/api"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/constantes"
	"github.com/Nhanderu/brdoc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	app := fiber.New()

	v := validator.New()
	var j jsoniter.API

	var w sync.WaitGroup

	w.Add(2)

	go iniciarValidate(v, &w)
	go iniciarJson(&j, &w)

	var pageSize uint8

	if ps := os.Getenv("DATABASE_PAGESIZE"); ps != "" {
		p, err := strconv.ParseInt(ps, 10, 8)

		if err != nil {
			log.Println("DATABASE_PAGESIZE not valid setting 50")
			pageSize = 50
		} else {
			pageSize = uint8(p)
		}
	} else {
		pageSize = 50
	}

	db, err := banco.NewMysqlConn(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		pageSize,
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
	var err error

	if err = v.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		cpf := fl.Field().String()

		if cpf == "" || len(cpf) != 11 {
			return false
		}

		return brdoc.IsCPF(cpf)
	}); err != nil {
		panic(err)
	}

	if err = v.RegisterValidation("cep", func(fl validator.FieldLevel) bool {
		cep := fl.Field().String()

		if cep == "" || len(cep) != 8 {
			return false
		}

		return brdoc.IsCEP(
			cep,
			brdoc.AC, brdoc.AL, brdoc.AM, brdoc.AP, brdoc.BA, brdoc.CE, brdoc.DF,
			brdoc.ES, brdoc.GO, brdoc.MA, brdoc.MG, brdoc.MS, brdoc.MT, brdoc.PA,
			brdoc.PB, brdoc.PE, brdoc.PI, brdoc.PR, brdoc.RJ, brdoc.RN, brdoc.RO,
			brdoc.RR, brdoc.RS, brdoc.SC, brdoc.SE, brdoc.SP, brdoc.TO,
		)
	}); err != nil {
		panic(err)
	}

	if err = v.RegisterValidation("uf", func(fl validator.FieldLevel) bool {
		uf := fl.Field().String()

		if uf == "" || len(uf) != 2 {
			return false
		}

		return strings.Contains(constantes.Estados, uf)
	}); err != nil {
		panic(err)
	}

	if err = v.RegisterValidation("unix", func(fl validator.FieldLevel) bool {
		unix := fl.Field().Int()

		return time.Now().Unix() < unix
	}); err != nil {
		panic(err)
	}

	w.Done()
}
