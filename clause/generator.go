package clause

import (
	"fmt"
	"strings"
)

// generators are used to generate sql clauses
var generators map[Type]func(...interface{}) (string, []interface{})

func init() {
	generators = make(map[Type]func(...interface{}) (string, []interface{}))
	generators[SELECT] = _select
	generators[INSERT] = _insert
	generators[WHERE] = _where
	generators[VALUES] = _values
	generators[ORDERBY] = _orderBy
	generators[LIMIT] = _limit

}

// SELECT $fields FROM $tableName
func _select(values ...interface{}) (sql string, sqlArgs []interface{}) {
	tableName := values[0]
	fields := values[1].([]string)

	if len(fields) == 1 && fields[0] == "*" { // select *
		sql = fmt.Sprintf("SELECT * FROM `%v`", tableName)
	} else {
		fields := strings.Join(values[1].([]string), ",")
		sql = fmt.Sprintf("SELECT %v FROM `%v`", fields, tableName)
	}

	return
}

// INSERT INTO $tableName ($fields)
func _insert(values ...interface{}) (sql string, sqlArgs []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	sql = fmt.Sprintf("INSERT INTO `%v` (%s)", tableName, fields)

	return
}

// WHERE $desc	args
func _where(values ...interface{}) (sql string, sqlArgs []interface{}) {
	desc := values[0] // like 'name = ? and age > ?'
	sql = fmt.Sprintf("WHERE %s", desc)
	sqlArgs = values[1:]

	return
}

// VALUES (?,?,?...),(?,?,?...),(?,?,?...)	args
func _values(values ...interface{}) (sql string, sqlArgs []interface{}) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("VALUES ")
	descs := make([]string, 0)
	for _, value := range values {
		v := value.([]interface{})
		n := len(v)
		elems := make([]string, n)
		for i := range elems {
			elems[i] = "?"
		}
		desc := strings.Join(elems, ",")
		descs = append(descs, fmt.Sprintf("(%v)", desc))
		sqlArgs = append(sqlArgs, v...)
	}
	sqlBuilder.WriteString(strings.Join(descs, ","))
	sql = sqlBuilder.String()

	return
}

// ORDERBY $desc
func _orderBy(values ...interface{}) (sql string, sqlArgs []interface{}) {
	order := "ASC"
	if len(values) > 1 {
		order = values[1].(string)
	}
	sql = fmt.Sprintf("ORDER BY %s %s", values[0], order)

	return
}

// LIMIT ?	args
func _limit(values ...interface{}) (sql string, sqlArgs []interface{}) {
	sql = "LIMIT ?"
	sqlArgs = values

	return
}
