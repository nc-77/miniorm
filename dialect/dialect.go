package dialect

import "reflect"

var dialectMap = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSql(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
