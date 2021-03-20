package api

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strconv"
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

	apiColaborador := app.Group("/api/v1/colaborador", func(ctx *fiber.Ctx) error {
		if token := ctx.Get(TOKEN_HEADER, "a"); token != "a" {
			if email := ctx.Get(EMAIL_HEADER, "a"); email != "a" {
				if valid, _ := db.IsValidToken(email, token); valid {
					return ctx.Next()
				}
			}
		}

		return ctx.SendStatus(http.StatusForbidden)
	})

	apiColaborador.Get("/produtos/listar/:page", func(ctx *fiber.Ctx) error {
		page, err := strconv.Atoi(ctx.Params("page"))

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := ctx.Get(TOKEN_HEADER, "a")
		if token == "a" {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarProdutos(token, uint8(page)))
	})

	apiAdministrador := app.Group("/api/v1/administrador", func(ctx *fiber.Ctx) error {
		if token := ctx.Get(TOKEN_HEADER, "a"); token != "a" {
			if email := ctx.Get(EMAIL_HEADER, "a"); email != "a" {
				if valid, admin := db.IsValidToken(email, token); valid && admin {
					return ctx.Next()
				}
			}
		}

		return ctx.SendStatus(http.StatusForbidden)
	})

	apiAdministrador.Put("/usuarios/cadastrar", func(ctx *fiber.Ctx) error {
		var r objetos.Usuario

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.CadastarUsuario(ctx.Get(TOKEN_HEADER), r); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrador.Get("/usuarios/listar/:page", func(ctx *fiber.Ctx) error {
		page, err := strconv.Atoi(ctx.Params("page"))

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := ctx.Get(TOKEN_HEADER, "a")
		if token == "a" {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarUsuarios(token, uint8(page)))
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
