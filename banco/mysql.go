package banco

import (
	"context"
	"database/sql"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type MysqlDriver struct {
	*sql.Conn
}

func NewMysqlConn(host, port, username, password, database string) (DriverBancoDados, error) {
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

	return &MysqlDriver{conn}, err
}

func (m MysqlDriver) CadastarUsuario(u objetos.Usuario) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarUsuario(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		u.UNOME, u.UEMAIL, u.UPASSWORD, u.UCPF, u.UUF,
		u.UCIDADE, u.UCEP, u.UENDERECO, u.UCOMPLEMENTO,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m MysqlDriver) Login(uemail, upassword string) string {
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

func (m MysqlDriver) IsValidToken(uemail, utoken string) (bool, string) {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ValidarToken(?, ?)",
		uemail, utoken,
	)

	if err != nil {
		return false,""
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

func (m MysqlDriver) Logoff(uemail, token string) uint8 {
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

func (m MysqlDriver) ListarEspecialidades() []objetos.Especialidade {
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

func (m MysqlDriver) ListarMedicoPorEspecialiade(eid uint) []objetos.Medico {
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

func (m MysqlDriver) ListarAgendamentosDoMedico(mid uint64) []objetos.Agendamento {
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

func (m MysqlDriver) ListarHospitais() []objetos.Hospital {
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

func (m MysqlDriver) MarcarConsulta(utoken string, mid uint64, data time.Time) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarAgendamento(?, ?, ?)", //TODO GARANTIR LIMITE DE MARCAÇÃO DE CONSULTA
		utoken, mid, data,
	)

	if err != nil {
		return 1
	}

	return 0
}

func (m MysqlDriver) AdicionarHospital(hospital objetos.Hospital) uint8 {
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

func (m MysqlDriver) AdicionarMedico(medico objetos.Medico) uint8 {
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

func (m MysqlDriver) AdicionarEspecialidade(nome string) uint8 {
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
