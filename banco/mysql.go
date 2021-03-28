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
	pageSize uint8
}

func NewMysqlConn(host, port, username, password, database string, pageSize uint8) (DriverBancoDados, error) {
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

	d, err := sql.Open("mysql", connString.String())

	if err != nil {
		return nil, err
	}

	conn, err := d.Conn(context.Background())

	db := &MysqlDriver{conn, pageSize}

	return db, err
}

func (m MysqlDriver) CadastarUsuario(u objetos.Usuario) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL AddUsuario(?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		return ""
	}

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

func (m MysqlDriver) IsValidToken(uemail, utoken string) bool {
	r, err := m.QueryContext(
		context.Background(),
		"CALL ValidarToken(?, ?)",
		uemail, utoken,
	)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		return false
	}

	if r.Next() {
		var i uint8
		if err = r.Scan(&i); err != nil {
			return false
		} else {
			return i == 1
		}
	}

	return false
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

func (m MysqlDriver) ListarEspecialidades(page uint8) []objetos.Especialidade {
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
	i := 0

	for r.Next() {
		e := objetos.Especialidade{}

		if err = r.Scan(&e.EID, &e.ENome); err != nil {
			return list
		}

		list = append(list, e)

		i++
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
	i := 0

	for r.Next() {
		M := objetos.Medico{}

		if err = r.Scan(&M.MID, &M.HID, &M.MNOME); err != nil {
			return list
		}

		list = append(list, M)

		i++
	}

	return list
}

func (m MysqlDriver) ListarAgendamentosDoMedico(mid uint64, page uint8) []objetos.Agendamento {
	return make([]objetos.Agendamento, 0)
}

func (m MysqlDriver) MarcarConsulta(utoken string, mid uint64, data time.Time) uint8 {
	_, err := m.QueryContext(
		context.Background(),
		"CALL RegistrarAgendamento(?, ?, ?)",
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
