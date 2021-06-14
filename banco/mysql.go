//+build mysql-driver
package banco

import (
	"context"
	"database/sql"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	. "github.com/icza/gox/gox"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type Driver struct {
	*sql.Conn
}

func (m Driver) NewConn(host, port, username, password, database string) (DriverBancoDados, error) {
	var connString strings.Builder

	//username:password@tcp(host:port)/database
	connString.WriteString(username)
	connString.WriteString(":")
	connString.WriteString(password)
	connString.WriteString("@tcp(")
	connString.WriteString(host)
	connString.WriteString(":")
	connString.WriteString(port)
	connString.WriteString(")/")
	connString.WriteString(database)
	connString.WriteString("?parseTime=true")

	d, err := sql.Open("mysql", connString.String())

	if err != nil {
		return nil, err
	}

	conn, err := d.Conn(context.Background())

	return &Driver{conn}, err
}

func (m Driver) CadastarUsuario(u objetos.Usuario) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarUsuario(?, ?, ?, ?, ?, ?, ?, ?)",
		u.UNOME, u.UEMAIL, u.UPASSWORD, u.UCPF, time.Unix(u.UNASCIMENTO, 0),
		u.USEXO, u.UTELEFONE, u.CPID,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) Login(uemail, upassword string) string {
	r, err := m.QueryContext(
		context.Background(),
		"CALL GetLoginUsuario(?)",
		uemail,
	)

	if err != nil {
		return ""
	}

	defer func() {
		_ = r.Close()
	}()

	if r.Next() {
		var password string
		err = r.Scan(&password)

		if err != nil || r.Close() != nil {
			return ""
		}

		if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(upassword)); err == nil {
			token := uuid.NewString()

			_, err = m.QueryContext(
				context.Background(),
				"CALL RegistrarToken(?, ?)",
				uemail, token,
			)

			if err != nil {
				return ""
			}

			return token
		}

		return ""
	}

	return ""
}

func (m Driver) IsValidToken(uemail, utoken string) (bool, string) {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ValidarToken(?, ?)",
		uemail, utoken,
	)

	if err != nil {
		return false, ""
	}

	defer func() {
		_ = r.Close()
	}()

	if r.Next() {
		var i uint8
		if err = r.Scan(&i); err != nil {
			return false, ""
		} else {
			if i == 1 {
				token := uuid.NewString()
				_ = r.Close()
				r, err = m.QueryContext(
					context.Background(),
					"CALL RegistrarToken(?, ?)",
					uemail, token,
				)

				return true, token
			}

			return false, ""
		}
	}

	return false, ""
}

func (m Driver) Logoff(uemail, token string) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL LogOff(? , ?)",
		uemail, token,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) ListarEspecialidades() []objetos.Especialidade {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarEspecialidades()",
	)

	if err != nil {
		return make([]objetos.Especialidade, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Especialidade

	for r.Next() {
		e := objetos.Especialidade{}

		if err = r.Scan(&e.EID, &e.ENome); err != nil {
			return list
		}

		list = append(list, e)
	}

	return list
}

func (m Driver) ListarEspecialidadesHospital(hid uint) []uint {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarEspecialidadesHospital(?)",
		hid,
	)

	if err != nil {
		return make([]uint, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []uint

	for r.Next() {
		var e uint

		if err = r.Scan(&e); err != nil {
			return list
		}

		list = append(list, e)
	}

	return list
}

func (m Driver) ListarAgendamentos(utoken string) []objetos.Agendamento {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarAgendamentosUsuario(?)",
		utoken,
	)

	if err != nil {
		return make([]objetos.Agendamento, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Agendamento

	for r.Next() {
		a := objetos.Agendamento{}

		var data time.Time

		if err = r.Scan(&a.AID, &data, &a.MID, &a.DID); err != nil {
			return list
		}

		a.ADATA = data.Unix()

		list = append(list, a)
	}

	return list
}

func (m Driver) FavoritarHospital(utoken string, hid uint) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL FavoritarHospital(?, ?)",
		utoken, hid,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) ListarHospitaisFavoritos(utoken string) []objetos.Hospital {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarHospitaisFavoritos(?)",
		utoken,
	)

	if err != nil {
		return make([]objetos.Hospital, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Hospital

	for r.Next() {
		h := objetos.Hospital{}

		if err = r.Scan(&h.HID, &h.HNOME, &h.HUF, &h.HCIDADE, &h.HCEP,
			&h.HENDERECO, &h.HCOMPLEMENTO, &h.HTELEFONE, &h.HISPRONTOSOCORRO,
		); err != nil {
			return list
		}

		list = append(list, h)
	}

	return list
}

func (m Driver) AdicionarDependete(utoken string, dependete objetos.Dependente) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarDependete(?, ?, ?, ?)",
		utoken, dependete.DNOME, time.Unix(dependete.DNASCIMENTO, 0),
		dependete.DSEXO,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) ListarDependentes(utoken string) []objetos.Dependente {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarDependentes(?)",
		utoken,
	)

	if err != nil {
		return make([]objetos.Dependente, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Dependente

	for r.Next() {
		d := objetos.Dependente{}
		var data time.Time

		if err = r.Scan(&d.DID, &d.DNOME, &data, &d.DSEXO); err != nil {
			return list
		}

		d.DNASCIMENTO = data.Unix()

		list = append(list, d)
	}

	return list
}

func (m Driver) RemoverDependente(utoken string, did uint64) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL DesativarDepente(?, ?)",
		utoken, did,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) ListarConvenios() []objetos.Convenio {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarConvenios()",
	)

	if err != nil {
		return make([]objetos.Convenio, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Convenio

	for r.Next() {
		c := objetos.Convenio{}

		if err = r.Scan(&c.CID, &c.CNOME); err != nil {
			return list
		}

		list = append(list, c)
	}

	return list
}

func (m Driver) ListarPlanosConvenio(cid uint) []objetos.Planos {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarPlanos(?)",
		cid,
	)

	if err != nil {
		return make([]objetos.Planos, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Planos

	for r.Next() {
		p := objetos.Planos{}

		if err = r.Scan(&p.CPID, &p.CPNOME, &p.CID); err != nil {
			return list
		}

		list = append(list, p)
	}

	return list
}

func (m Driver) AdicionarConvenioHospital(cpid uint64, hid uint) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarConvenioHospital(?, ?)",
		hid, cpid,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) RemoverConvenioHospital(cpid uint64, hid uint) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RemoverConvenioHospital(?, ?)",
		cpid, hid,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) AdicionarConvenio(nome string) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarConvenio(?)",
		nome,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) AdicionarPlanoConvenio(cid uint64, nome string) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistarPlanoEmConvenio(?, ?)",
		cid, nome,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) ListarHospitaisPorPlanoConvenio(cpid uint64) []uint {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarHospitaisPorPlanoConvenio(?)",
		cpid,
	)

	if err != nil {
		return make([]uint, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []uint

	for r.Next() {
		var hid uint

		if err = r.Scan(&hid); err != nil {
			return list
		}

		list = append(list, hid)
	}

	return list
}

func (m Driver) ListarMedicoPorEspecialiade(eid uint) []objetos.Medico {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarMedicosPorEspecialidade(?)",
		eid,
	)

	if err != nil {
		return make([]objetos.Medico, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Medico

	for r.Next() {
		M := objetos.Medico{}

		if err = r.Scan(&M.MID, &M.HID, &M.MNOME); err != nil {
			return list
		}

		list = append(list, M)
	}

	return list
}

func (m Driver) ListarAgendamentosDoMedico(mid uint64) []objetos.Agendamento {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarAgendamentosMedico(?)",
		mid,
	)

	if err != nil {
		return make([]objetos.Agendamento, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Agendamento

	for r.Next() {
		a := objetos.Agendamento{}

		if err = r.Scan(&a.ADATA); err != nil {
			return list
		}

		list = append(list, a)
	}

	return list
}

func (m Driver) ListarHospitais() []objetos.Hospital {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ListarHospitais()",
	)

	if err != nil {
		return make([]objetos.Hospital, 0)
	}

	defer func() {
		_ = r.Close()
	}()

	var list []objetos.Hospital

	for r.Next() {
		h := objetos.Hospital{}

		if err = r.Scan(&h.HID, &h.HNOME, &h.HUF, &h.HCIDADE, &h.HCEP,
			&h.HENDERECO, &h.HCOMPLEMENTO, &h.HTELEFONE, &h.HISPRONTOSOCORRO,
		); err != nil {
			return list
		}

		list = append(list, h)
	}

	return list
}

func (m Driver) MarcarConsulta(utoken string, did uint64, mid uint64, data time.Time) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarAgendamento(?, ?, ?, ?)", //TODO GARANTIR LIMITE DE MARCAÇÃO DE CONSULTA
		utoken, If(did == 0).If(nil, did), mid, data,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) AdicionarHospital(hospital objetos.Hospital) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarHospital(?, ?, ?, ?, ?, ?, ?, ?)",
		hospital.HNOME, hospital.HUF, hospital.HCIDADE, hospital.HCEP,
		hospital.HENDERECO, hospital.HCOMPLEMENTO, hospital.HTELEFONE, hospital.HISPRONTOSOCORRO,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) AdicionarMedico(medico objetos.Medico) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarMedico(?, ?, ?)",
		medico.HID, medico.EID, medico.MNOME,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m Driver) AdicionarEspecialidade(nome string) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarEspecialidade(?)",
		nome,
	)

	if err != nil {
		return 1
	}

	return 0
}
