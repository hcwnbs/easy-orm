package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Table interface {
	TableName() string
}

// EngineForMysql EngineForMysql for mysql
type EngineForMysql struct {
	// original sql db
	Db *sql.DB
	Tx *sql.Tx

	TableName string
	Prepare   string
	Exec      []interface{}
}

func NewMySqlEngine(driverName, useName, password, address, dbName string) (*EngineForMysql, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&timeout=5s&readTimeout=6s", useName, password, address, dbName)
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	//Todo 最大连接数等配置

	return &EngineForMysql{Db: db}, nil
}


// Table 设置表名
func (e *EngineForMysql) Table(tableName string) *EngineForMysql {
	e.TableName = tableName
	return e
}

// GetTable 获取表名
func (e *EngineForMysql) GetTable() string {
	return e.TableName
}

func (e *EngineForMysql) Insert(data interface{}) error {
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

func (e *EngineForMysql) exec() error {
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

func (e *EngineForMysql) resetEngine() {
	e.TableName = ""
	e.Prepare = ""
	e.Exec = []interface{}{}
}
