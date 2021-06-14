package api

import (
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/banco"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

const (
	TOKEN_HEADER = "token_api"
	EMAIL_HEADER = "email_api"
)

func BuildRoutes(app *fiber.App, db banco.DriverBancoDados, v *validator.Validate, json *jsoniter.API) {
	app.Get("/api/v1/convenios", func(ctx *fiber.Ctx) error {
		return ctx.JSON(db.ListarConvenios())
	})

	app.Get("/api/v1/convenio/:cid/planos", func(ctx *fiber.Ctx) error {
		cid, err := strconv.ParseUint(ctx.Params("cid"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if cid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarPlanosConvenio(uint(cid)))
	})

	app.Post("/api/v1/cadastrar/usuario", func(ctx *fiber.Ctx) error {
		var r objetos.Usuario

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(r.UPASSWORD), bcrypt.DefaultCost)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		r.UPASSWORD = string(pass)

		if e := db.CadastarUsuario(r); e != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

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
			Token string `json:"token" validate:"uuid_rfc4122,required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		db.Logoff(r.Email, r.Token) //TODO STATUS ERROR MAP

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario := app.Group("/api/v1/usuario", func(ctx *fiber.Ctx) error {
		if token := ctx.Get(TOKEN_HEADER, "a"); token != "a" {
			if email := ctx.Get(EMAIL_HEADER, "a"); email != "a" {
				if err := v.Var(email, "email"); err == nil {
					if err = v.Var(token, "uuid_rfc4122"); err == nil {
						if valido, t := db.IsValidToken(email, token); valido {
							ctx.Set(TOKEN_HEADER, t)
							return ctx.Next()
						}
					}
				}
			}
		}

		return ctx.SendStatus(http.StatusForbidden)
	})

	apiUsuario.Get("/especialidades", func(ctx *fiber.Ctx) error {
		return ctx.JSON(db.ListarEspecialidades())
	})

	apiUsuario.Get("/hospitais", func(ctx *fiber.Ctx) error {
		return ctx.JSON(db.ListarHospitais())
	})

	apiUsuario.Post("/favoritos/hospital/:hid/add", func(ctx *fiber.Ctx) error {
		hid, err := strconv.ParseUint(ctx.Params("hid"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if hid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		if err := db.FavoritarHospital(string(token), uint(hid)); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario.Get("/favoritos/hospital", func(ctx *fiber.Ctx) error {
		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		return ctx.JSON(db.ListarHospitaisFavoritos(string(token)))
	})

	apiUsuario.Get("/dependetes", func(ctx *fiber.Ctx) error {
		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		return ctx.JSON(db.ListarDependentes(string(token)))
	})

	apiUsuario.Put("/dependete/add", func(ctx *fiber.Ctx) error {
		var r objetos.Dependente

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		if err := db.AdicionarDependete(string(token), r); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario.Delete("/dependete/:did/del", func(ctx *fiber.Ctx) error {
		did, err := strconv.ParseUint(ctx.Params("did"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if did < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		if err := db.RemoverDependente(string(token), did); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario.Get("/hospital/:hid/especialidades", func(ctx *fiber.Ctx) error {
		hid, err := strconv.ParseUint(ctx.Params("hid"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if hid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarEspecialidadesHospital(uint(hid)))
	})

	apiUsuario.Get("/hospital/:eid/medicos", func(ctx *fiber.Ctx) error {
		eid, err := strconv.ParseUint(ctx.Params("eid"), 10, strconv.IntSize)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if eid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarMedicoPorEspecialiade(uint(eid)))
	})

	apiUsuario.Get("/hospital/:medico/agenda", func(ctx *fiber.Ctx) error {
		mid, err := strconv.ParseUint(ctx.Params("medico"), 10, 64)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if mid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarAgendamentosDoMedico(mid))
	})

	apiUsuario.Put("/hospital/:medico/agenda/add", func(ctx *fiber.Ctx) error {
		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		mid, err := strconv.ParseUint(ctx.Params("medico"), 10, 64)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if mid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		var r struct {
			Data int64  `json:"data" validate:"unix-futuro,required"`
			Did  uint64 `json:"did"`
		}

		if err = getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.MarcarConsulta(string(token), r.Did, mid, time.Unix(r.Data, 0)); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiUsuario.Get("/agenda", func(ctx *fiber.Ctx) error {
		token := ctx.Response().Header.Peek(TOKEN_HEADER)
		return ctx.JSON(db.ListarAgendamentos(string(token)))
	})

	apiUsuario.Get("/convenio/:cpid/hospitais", func(ctx *fiber.Ctx) error {
		cpid, err := strconv.ParseUint(ctx.Params("cpid"), 10, 64)

		if err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if cpid < 0 {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.JSON(db.ListarHospitaisPorPlanoConvenio(cpid))
	})

	apiAdministrativa := app.Group("/api/v1/adm", func(ctx *fiber.Ctx) error {
		return ctx.Next() //TODO++ SUPER AUTENTICAÇÃO
	})

	apiAdministrativa.Put("/hospital/add", func(ctx *fiber.Ctx) error {
		var r objetos.Hospital

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarHospital(r); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Put("/hospital/convenio/add", func(ctx *fiber.Ctx) error {
		var r struct {
			CPID uint64 `json:"cpid" validate:"required"`
			HID  uint   `json:"hid" validate:"required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarConvenioHospital(r.CPID, r.HID); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Delete("/hospital/convenio/del", func(ctx *fiber.Ctx) error {
		var r struct {
			CPID uint64 `json:"cpid" validate:"required"`
			HID  uint   `json:"hid" validate:"required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.RemoverConvenioHospital(r.CPID, r.HID); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Put("/hospital/medico/add", func(ctx *fiber.Ctx) error {
		var r objetos.Medico

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarMedico(r); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Put("/especialidades/add", func(ctx *fiber.Ctx) error {
		var r struct {
			Nome string `json:"nome"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarEspecialidade(r.Nome); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Put("/convenio/add", func(ctx *fiber.Ctx) error {
		var r struct {
			Nome string `json:"nome"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarConvenio(r.Nome); err != 0 {
			return ctx.SendStatus(http.StatusConflict)
		}

		return ctx.SendStatus(http.StatusOK)
	})

	apiAdministrativa.Put("/convenio/plano/add", func(ctx *fiber.Ctx) error {
		var r struct {
			Nome string `json:"nome" validate:"required"`
			CID  uint64 `json:"cid" validate:"required"`
		}

		if err := getRequest(v, *json, ctx, &r); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		if err := db.AdicionarPlanoConvenio(r.CID, r.Nome); err != 0 {
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
