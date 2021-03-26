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

	d, err := sql.Open("mysql", connString.String())

	if err != nil {
		return nil, err
	}

	conn, err := d.Conn(context.Background())

	db := &MysqlDriver{conn}

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

func (m MysqlDriver) ListarEspecialidades(page uint8) map[uint]string {
	return make(map[uint]string)
}

func (m MysqlDriver) ListarMedicoPorEspecialiade(eid uint) []objetos.Medico {
	return make([]objetos.Medico, 0)
}

func (m MysqlDriver) ListarAgendamentosDoMedico(mid uint64, page uint8) map[uint64]time.Time {
	return make(map[uint64]time.Time)
}

func (m MysqlDriver) MarcarConsulta(utoken string, mid uint64, data time.Time) uint8 {
	return 0
}
