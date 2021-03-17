package banco

import (
	"database/sql"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type MysqlDriver struct {
	*sql.DB
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

	db := &MysqlDriver{d}

	return db, err
}

func (m MysqlDriver) Login(lid uint, uemail, upassword string) string {
	panic("implement me")
}

func (m MysqlDriver) Logoff(uemail, token string) uint {
	panic("implement me")
}

func (m MysqlDriver) CadastarUsuario(lid uint, tokenAdmin string, usuario objetos.Usuario) uint {
	panic("implement me")
}

func (m MysqlDriver) AtualizarUsuario(lid uint, tokenAdmin string, usuario objetos.Usuario) uint {
	panic("implement me")
}

func (m MysqlDriver) ListarUsuarios(lid uint, tokenAdmin string) []objetos.Usuario {
	panic("implement me")
}

func (m MysqlDriver) CriarProduto(lid uint, token string, produto objetos.Produto) uint {
	panic("implement me")
}

func (m MysqlDriver) ListarProdutos(lid uint, token string, page uint8) []objetos.Produto {
	panic("implement me")
}

func (m MysqlDriver) AtualizarProduto(lid uint, token string, produto objetos.Produto) uint {
	panic("implement me")
}

func (m MysqlDriver) RegistrarOrdem(lid uint, token string, ordem objetos.Ordem) uint {
	panic("implement me")
}

func (m MysqlDriver) ListarEstoque(lid uint, token string, page uint8) []objetos.ItemEstoque {
	panic("implement me")
}
