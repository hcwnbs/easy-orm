package easy_orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const MysqlDriver = "mysql"

type Table interface {
	TableName() string
}

// EasyOrmEngine core engine
type EasyOrmEngine struct {
	// original sql db
	Db *sql.DB
	Tx *sql.Tx

	TableName string
	Prepare   string
	Exec      []interface{}
}

func NewEngine(driverName, useName, password, address, dbName string) (*EasyOrmEngine, error) {
	switch driverName {
	case MysqlDriver:
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=5s&readTimeout=6s", useName, password, address, dbName)
		db, err := sql.Open(MysqlDriver, dsn)
		if err != nil {
			return nil, err
			}

		//Todo 最大连接数等配置

		return &EasyOrmEngine{Db: db}, nil
	default:
		return nil, fmt.Errorf("the driver %s is not support", driverName)
	}
}

// Table 设置表名
func (e *EasyOrmEngine) Table(tableName string) *EasyOrmEngine {
	e.TableName = tableName
	return e
}

// GetTable 获取表名
func (e *EasyOrmEngine) GetTable() string {
	return e.TableName
}

func (e *EasyOrmEngine) Insert(data interface{}) error {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("what to inster is not a struct")
	}

	// 字段名
	var fieldName []string

	// ?占位符
	var placeholder []string

	for i := 0; i < t.NumField(); i++ {
		//如果是字段是小写开头则无法反射
		if !v.Field(i).CanInterface() {
			continue
		}

		tag := t.Field(i).Tag.Get("orm")
		if tag != "" {
			// 跳过自增字段
			if strings.Contains(tag, "auto_increment") {
				continue
			}
			fieldName = append(fieldName, tag)
		} else {
			fieldName = append(fieldName, strings.ToLower(t.Field(i).Name))
		}

		placeholder = append(placeholder, "?")
		e.Exec = append(e.Exec, v.Field(i).Interface())
	}

	// 如果未设置表名
	if e.TableName == "" {
		if table, ok := data.(Table); ok {
			e.Table(table.TableName())
		} else {
			e.Table(strings.ToLower(t.Name()))
		}
	}

	e.Prepare = fmt.Sprintf("insert into %s (%s) values (%s)",
		e.GetTable(), strings.Join(fieldName, ","), strings.Join(placeholder, ","))

	return e.exec()
}

func (e *EasyOrmEngine) exec() error {
	stmt, err := e.Db.Prepare(e.Prepare)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Exec...)
	if err != nil {
		return err
	}
	e.resetEngine()
	return nil
}

func (e *EasyOrmEngine) resetEngine() {
	e.TableName = ""
	e.Prepare = ""
	e.Exec = []interface{}{}
}




