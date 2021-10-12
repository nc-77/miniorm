package session

import (
	"database/sql"
	"errors"
	"miniorm/clause"
	"reflect"
)

var (
	ErrRecordNotFound = errors.New("record not find")
)

// FindRecords assign all records to the values by parsing values
func (s *Session) FindRecords(values interface{}) (err error) {
	var affected int64
	defer func() {
		s.recordLast(affected, err)
	}()

	// check values kind
	distSlice := reflect.Indirect(reflect.ValueOf(values))
	distType := distSlice.Type().Elem()

	// build sql and sqlArgs based on refTable
	table := s.RefTable()
	s.clause.Set(clause.SELECT, table.Name, []string{"*"})
	sqls, sqlArgs := s.clause.Build(clause.SELECT)

	// query rows and map rows to values
	var rows *sql.Rows
	var columns []string
	if rows, err = s.Raw(sqls, sqlArgs...).Query(); err != nil {
		return
	}
	if columns, err = rows.Columns(); err != nil {
		return
	}
	for rows.Next() {
		dist := reflect.New(distType).Elem()
		scanValues := make([]interface{}, len(columns))
		for i := range scanValues {
			scanValues[i] = dist.FieldByName(columns[i]).Addr().Interface()
		}
		if err = rows.Scan(scanValues...); err != nil {
			return
		}
		distSlice.Set(reflect.Append(distSlice, dist))
		affected++
	}
	return rows.Close()
}

// FirstRecord assign the first record to value
func (s *Session) FirstRecord(value interface{}) (err error) {
	//todo
	return
}

// CreateRecord insert the value into database by parsing value
func (s *Session) CreateRecord(value interface{}) (err error) {
	// todo
	return
}
