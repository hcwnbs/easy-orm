package engine

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"easy-orm/engine/mysql"
)

const MysqlDriver = "mysql"

type Engine interface {
    Insert(data interface{}) error
}

func NewEngine(driverName, useName, password, address, dbName string) (Engine, error) {
	switch driverName {
	case MysqlDriver:
		return mysql.NewMySqlEngine(driverName, useName, password, address, dbName)
	default:
		return nil, fmt.Errorf("the driver is not support, driver: %s", driverName)
	}
}
