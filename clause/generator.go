package clause

import (
	"fmt"
	"reflect"
	"strings"
)

// generators are used to generate sql clauses
var generators map[Type]func(...interface{}) (string, []interface{})

func init() {
	generators = make(map[Type]func(...interface{}) (string, []interface{}))
	generators[SELECT] = _select
	generators[INSERT] = _insert
	generators[DELETE] = _delete
	generators[UPDATE] = _update
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

// DELETE FROM $tableName
func _delete(values ...interface{}) (sql string, sqlArgs []interface{}) {
	tableName := values[0]
	sql = fmt.Sprintf("DELETE FROM %v", tableName)

	return
}

// UPDATE $tableName SET $field1=$value1,$field2=$value2
func _update(values ...interface{}) (sql string, sqlArgs []interface{}) {
	tableName := values[0]
	value := reflect.Indirect(reflect.ValueOf(values[1]))
	fields := make([]string, 0)

	switch value.Kind() {
	case reflect.Map:
		fieldVals := values[1].(map[string]interface{})
		for k, v := range fieldVals {
			fields = append(fields, fmt.Sprintf("%v = ?", k))
			sqlArgs = append(sqlArgs, v)
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fields = append(fields, fmt.Sprintf("%v = ?", value.Type().Field(i).Name))
			sqlArgs = append(sqlArgs, field.Interface())
		}
	}

	sql = fmt.Sprintf("UPDATE `%v` SET %v", tableName, strings.Join(fields, ","))

	return
}

// WHERE $desc	args
func _where(values ...interface{}) (sql string, sqlArgs []interface{}) {
	desc := values[0] // like 'name = ? and age > ?'
	sql = fmt.Sprintf("WHERE %v", desc)
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
