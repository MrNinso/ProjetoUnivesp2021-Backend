package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type MysqlDriver struct {
	*sql.DB
}

func NewMysqlConn(host, port, username, password, database string) (Database, error) {
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

func (m MysqlDriver) login() {
	panic("implement me")
}
