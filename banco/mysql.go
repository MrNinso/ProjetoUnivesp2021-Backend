package banco

import (
	"context"
	"database/sql"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	_ "github.com/go-sql-driver/mysql"
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

func (m MysqlDriver) CadastarUsuario(usuario objetos.Usuario) uint8 {
	return 0
}

func (m MysqlDriver) Login(uemail, upassword string) string {
	return "token"
}

func (m MysqlDriver) IsValidToken(uemail, utoken string) bool {
	return true
}

func (m MysqlDriver) Logoff(uemail, token string) uint8 {
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
	panic("implement me")
}
