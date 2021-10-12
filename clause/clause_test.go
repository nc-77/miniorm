package clause

import (
	"reflect"
	"testing"
)

func TestClause(t *testing.T) {
	var c Clause
	assert := func(sql string, sqlArgs []interface{}, sqlExcepted string, sqlArgsExcepted []interface{}) {
		if sql != sqlExcepted {
			t.Fatalf("build sql failed\nexcepted: %v\ngot     : %v", sqlExcepted, sql)
		}
		for i := range sqlArgs {
			if !reflect.DeepEqual(sqlArgs[i], sqlArgsExcepted[i]) {
				t.Fatalf("build sqlArgs failed\nexcepted: %v\ngot     : %v", sqlArgsExcepted[i], sqlArgs[i])
			}
		}
	}

	t.Run("SELECT_0", func(t *testing.T) {
		c.Set(SELECT, "User", []string{"*"})
		c.Set(WHERE, "name = ?", "nic")
		c.Set(ORDERBY, "id")
		c.Set(LIMIT, 1)
		sql, sqlArgs := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
		sqlExcepted := "SELECT * FROM `User` WHERE name = ? ORDER BY id ASC LIMIT ?"
		sqlArgsExcepted := []interface{}{"nic", 1}
		assert(sql, sqlArgs, sqlExcepted, sqlArgsExcepted)
	})

	t.Run("SELECT_1", func(t *testing.T) {
		c.Set(SELECT, "User", []string{"id", "name", "age"})
		c.Set(WHERE, "name = ? and age > ?", "nic", 18)
		c.Set(ORDERBY, "id", "DESC")
		c.Set(LIMIT, 1)
		sql, sqlArgs := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
		sqlExcepted := "SELECT (id,name,age) FROM `User` WHERE name = ? and age > ? ORDER BY id DESC LIMIT ?"
		sqlArgsExcepted := []interface{}{"nic", 18, 1}
		assert(sql, sqlArgs, sqlExcepted, sqlArgsExcepted)
	})

	//t.Run("INSERT_0", func(t *testing.T) {
	//	c.Set(INSERT,"User",[]string{})
	//})
}
