package clause

import "strings"

type Clause struct {
	sql     map[Type]string
	sqlArgs map[Type][]interface{}
}

type Type int

const (
	SELECT Type = iota
	INSERT
	UPDATE
	DELETE
	WHERE
	VALUES
	ORDERBY
	LIMIT
	SET
)

// Clear clears Clause's sql and sqlArgs
func (c *Clause) Clear() {
	c.sql = nil
	c.sqlArgs = nil
}

// Set sets Clause's sql and sqlArgs
func (c *Clause) Set(t Type, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlArgs = make(map[Type][]interface{})
	}
	sql, sqlArgs := generators[t](values...)
	c.sql[t], c.sqlArgs[t] = sql, sqlArgs
}

// Build combines the clauses into a complete sql and sqlArgs by orders
func (c *Clause) Build(orders ...Type) (sql string, sqlArgs []interface{}) {
	sqls := make([]string, 0)
	sqlArgs = make([]interface{}, 0)
	for _, order := range orders {
		sqls = append(sqls, c.sql[order])
		sqlArgs = append(sqlArgs, c.sqlArgs[order]...)
	}
	sql = strings.Join(sqls, " ")

	return
}

func (c *Clause) Exist(p Type) bool {
	_, ok := c.sql[p]
	return ok
}
