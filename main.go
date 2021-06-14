package main

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/api"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/constantes"
	"github.com/Nhanderu/brdoc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	. "github.com/icza/gox/gox"
	"github.com/icza/gox/timex"
	jsoniter "github.com/json-iterator/go"
	"os"
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

	db, err := banco.Driver{}.NewConn(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	Pie(err)

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

	Pie(app.Listen(port.String()))
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
	Pie(v.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		cpf := fl.Field().String()

		if cpf == "" || len(cpf) != 11 {
			return false
		}

		return brdoc.IsCPF(cpf)
	}))

	Pie(v.RegisterValidation("cep", func(fl validator.FieldLevel) bool {
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
	}))

	Pie(v.RegisterValidation("uf", func(fl validator.FieldLevel) bool {
		uf := fl.Field().String()

		if uf == "" || len(uf) != 2 {
			return false
		}

		return strings.Contains(constantes.Estados, uf)
	}))

	Pie(v.RegisterValidation("unix-futuro", func(fl validator.FieldLevel) bool {
		unix := fl.Field().Int()

		return time.Now().Unix() < unix //TODO ADD MARGEM
	}))

	Pie(v.RegisterValidation("unix-passado", func(fl validator.FieldLevel) bool {
		unix := fl.Field().Int()

		return time.Now().Unix() > unix //TODO ADD MARGEM
	}))

	Pie(v.RegisterValidation("idade", func(fl validator.FieldLevel) bool {
		unix := time.Unix(fl.Field().Int(), 0)

		idade, _, _, _, _, _ := timex.Diff(unix, time.Now())

		return idade >= 18
	}))

	Pie(v.RegisterValidation("sexo", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()

		return s == "M" || s == "F"
	}))

	w.Done()
}
