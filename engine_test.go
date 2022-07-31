package easy_orm

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"easy-orm/tests/prepare"
)

var TestMysqlEngine *EasyOrmEngine

func TestEasyOrmEngine_Insert(t *testing.T) {
	var user = prepare.User{
		Name: "test",
		Age:  1,
	}
	err := TestMysqlEngine.Insert(&user)
	require.Nil(t, err)

	err = TestMysqlEngine.Insert(user)
	require.Nil(t, err)

	err = TestMysqlEngine.Insert("test")
	require.Error(t, err)
}

// 测试原生的sql
func TestOriginalSql(t *testing.T) {
	// 插入测试数据
	stmt, err := TestMysqlEngine.Db.Prepare(`insert into user (name, age) values (?, ?)`)
	require.Nil(t, err)
	_, err = stmt.Exec(prepare.TestDefaultUserName, prepare.TestDefaultUserAge)
	require.Nil(t, err)

	// 单行查询测试数据
	var id1, age1 int
	var name1 string
	err = TestMysqlEngine.Db.QueryRow("select id, name, age from user where name = ?", prepare.TestDefaultUserName).Scan(
		&id1, &name1, &age1)
	require.Nil(t, err)
	require.Equal(t, prepare.TestDefaultUserName, name1)
	require.Equal(t, prepare.TestDefaultUserAge, age1)

	// 多行查询测试数据
	rows ,err := TestMysqlEngine.Db.Query("select id, name, age from user where name = ?", prepare.TestDefaultUserName)
	require.Nil(t, err)

	var users []prepare.User
	for rows.Next() {
		var id2, age2 int
		var name2 string
		err = rows.Scan(&id2, &name2, &age2)
		require.Nil(t, err)

		users = append(users, prepare.User{
			Id:   id2,
			Name: name2,
			Age:  age2,
		})
	}
	require.Equal(t, len(users), 1)
	if len(users) > 0 {
		require.Equal(t, prepare.TestDefaultUserName, users[0].Name)
		require.Equal(t, prepare.TestDefaultUserAge, users[0].Age)
	}

	//修改测试数据
	stmt, err = TestMysqlEngine.Db.Prepare(`update user set age = ? where name = ?`)
	require.Nil(t, err)
	_, err = stmt.Exec(prepare.TestDefaultUserAge + 1, prepare.TestDefaultUserName)
	require.Nil(t, err)

	//查询是否修改成功
	var age3 int
	err = TestMysqlEngine.Db.QueryRow("select age from user where name = ?", prepare.TestDefaultUserName).Scan(&age3)
	require.Nil(t, err)
	require.Equal(t, prepare.TestDefaultUserAge + 1, age3)

	//删除测试数据
	stmt, err = TestMysqlEngine.Db.Prepare(`delete from user where name = ?`)
	require.Nil(t, err)
	_, err = stmt.Exec(prepare.TestDefaultUserName)
	require.Nil(t, err)

	//查询是否删除成功
	var name3 string
	err = TestMysqlEngine.Db.QueryRow("select age from user where name = ?", prepare.TestDefaultUserName).Scan(&name3)
	require.Error(t, err)
}

func TestMain(m *testing.M) {
	engine, err := NewEngine(MysqlDriver, "root", "hcwnbs", "localhost:3306", "ut-test")
	if err != nil {
		fmt.Println("prepare for mysql engine failed, error:", err)
		os.Exit(1)
	}
	TestMysqlEngine = engine
	err = clearAll()
	if err != nil {
		fmt.Println("clear table failed, error: ", err)
		os.Exit(1)
	}
	m.Run()
	err = clearAll()
	if err != nil {
		fmt.Println("clear table failed, error: ", err)
		os.Exit(1)
	}
}

func clearAll() error {
	clearTables := []string{"user"}
	for _, table := range clearTables {
		sql := fmt.Sprintf("delete from %s where 1=1", table)
		_, err := TestMysqlEngine.Db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}