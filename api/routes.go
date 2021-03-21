package api

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strconv"
	"time"
)

const (
	TOKEN_HEADER = "token_api"
	EMAIL_HEADER = "email_api"
)

func BuildRoutes(app *fiber.App, db banco.DriverBancoDados, v *validator.Validate, json *jsoniter.API) {
	app.Post("/api/v1/login", func(ctx *fiber.Ctx) error {
		var r struct {
			Email    string `json:"email" validate:"email,required"`
			Password string `json:"password" validate:"required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if token := db.Login(r.Email, r.Password); token != "" {
			return ctx.JSON(fiber.Map{
				"token": token,
			})
		}

		return ctx.SendStatus(http.StatusForbidden)
	})

	app.Post("/api/v1/logoff", func(ctx *fiber.Ctx) error {
		var r struct {
			Email string `json:"email" validate:"email,required"`
			Token string `json:"token" validate:"alphanum,required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		db.Logoff(r.Email, r.Token) //TODO STATUS ERROR MAP

		return ctx.SendStatus(http.StatusOK)
	})

	app.Post("/api/v1/cadastrar/usuario", func(ctx *fiber.Ctx) error {
		var r objetos.Usuario

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.CadastarUsuario(r); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario := app.Group("/api/v1/usuario", func(ctx *fiber.Ctx) error {
		if token := ctx.Get(TOKEN_HEADER, "a"); token != "a" {
			if email := ctx.Get(EMAIL_HEADER, "a"); email != "a" {
				if db.IsValidToken(email, token) {
					return ctx.Next()
				}
			}
		}

		return ctx.SendStatus(http.StatusForbidden)
	})

	apiUsuario.Get("/especialiadades/:page", func(ctx *fiber.Ctx) error {
		p, err := strconv.ParseInt(ctx.Params("page"), 10, 8)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if p < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarEspecialidades(uint8(p)))
	})

	apiUsuario.Get("/hospital/medico/:eid", func(ctx *fiber.Ctx) error {
		eid, err := strconv.ParseUint(ctx.Params("eid"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if eid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarMedicoPorEspecialiade(uint(eid)))
	})

	apiUsuario.Get("/hospital/agenda/:medico/listar/:page", func(ctx *fiber.Ctx) error {
		mid, err := strconv.ParseUint(ctx.Params("medico"), 10, 64)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if mid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		p, err := strconv.ParseInt(ctx.Params("page"), 10, 8)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if p < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarAgendamentosDoMedico(mid, uint8(p)))
	})

	apiUsuario.Put("/hospital/agenda/:medico/add", func(ctx *fiber.Ctx) error {
		token := ctx.Get(TOKEN_HEADER)
		mid, err := strconv.ParseUint(ctx.Params("medico"), 10, 64)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if mid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		var r struct {
			data time.Time
		}

		if err = getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.MarcarConsulta(token, mid, r.data); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})
}

func getRequest(v *validator.Validate, j jsoniter.API, ctx *fiber.Ctx, r interface{}) error {
	if err := j.Unmarshal(ctx.Body(), r); err != nil {
		return err
	}

	if err := v.Struct(r); err != nil {
		return err
	}

	return nil
}
