package banco

import (
	"context"
	"database/sql"
	"github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"
	_ "github.com/go-sql-driver/mysql"
	"strings"
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

func (m MysqlDriver) Login(lid uint, uemail, upassword string) string {
	return "token" //TODO
}

func (m MysqlDriver) IsValidToken(uemail, utoken string) (isValid, isAdmin bool) {
	return true, true //TODO
}

func (m MysqlDriver) Logoff(lid uint, uemail, token string) uint {
	return 0 //TODO
}

func (m MysqlDriver) CadastarUsuario(token string, usuario objetos.Usuario) uint {
	return 0 //TODO
}

func (m MysqlDriver) AtualizarUsuario(token string, usuario objetos.Usuario) uint {
	return 0 //TODO
}

func (m MysqlDriver) ListarUsuarios(token string, page uint8) []objetos.Usuario {
	return make([]objetos.Usuario, 0) //TODO
}

func (m MysqlDriver) CriarProduto(token string, produto objetos.Produto) uint {
	return 0 //TODO
}

func (m MysqlDriver) ListarProdutos(token string, page uint8) []objetos.Produto {
	return make([]objetos.Produto, 0) //TODO
}

func (m MysqlDriver) AtualizarProduto(token string, produto objetos.Produto) uint {
	return 0 //TODO
}

func (m MysqlDriver) RegistrarOrdem(token string, ordem objetos.Ordem) uint {
	return 0 //TODO
}

func (m MysqlDriver) ListarEstoque(token string, page uint8) []objetos.ItemEstoque {
	return make([]objetos.ItemEstoque, 0) //TODO
}
