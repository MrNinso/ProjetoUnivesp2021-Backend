package api

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/gofiber/fiber/v2"
)

func BuildRoutes(app *fiber.App, db *banco.DriverBancoDados) {
	api := app.Group("/api/v1", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/usuario/login", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/usuario/logoff", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Put("/usuario/cadastrar", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/usuario/atualizar", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Get("/usuario/listar/:page", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/produto/cadastrar", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/produto/atualizar", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Get("/produto/listar/:page", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Post("/ordem/cadastrar", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Get("/ordem/listar/:page", func(ctx *fiber.Ctx) error {
		return nil
	})

	api.Get("/estoque/listar/:page", func(ctx *fiber.Ctx) error {
		return nil
	})
}
