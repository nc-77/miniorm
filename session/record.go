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
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	// build sql and sqlArgs based on refTable
	table := s.RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldsName)
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
		dest := reflect.New(destType).Elem()
		scanValues := make([]interface{}, len(columns))
		for i := range scanValues {
			scanValues[i] = dest.FieldByName(columns[i]).Addr().Interface()
		}
		if err = rows.Scan(scanValues...); err != nil {
			return
		}
		destSlice.Set(reflect.Append(destSlice, dest))
		affected++
	}
	err = rows.Close()

	return
}

// FirstRecord assign the first record to value
func (s *Session) FirstRecord(value interface{}) (err error) {
	return err
}

// CreateRecord insert the values into database by parsing value
func (s *Session) CreateRecords(values ...interface{}) (err error) {
	var result sql.Result
	defer func() {
		var affected int64
		if result != nil {
			affected, _ = result.RowsAffected()
		}
		s.recordLast(affected, err)
	}()

	table := s.RefTable()
	recordsValues := make([]interface{}, len(values))
	for i, value := range values {
		// map value to recordValues[i]
		destValue := reflect.Indirect(reflect.ValueOf(value))
		recordValues := make([]interface{}, len(table.FieldsName))
		for i, name := range table.FieldsName {
			recordValues[i] = destValue.FieldByName(name).Interface()
		}
		recordsValues[i] = recordValues
	}
	// build sql and sqlArgs based on refTable
	s.clause.Set(clause.INSERT, table.Name, table.FieldsName)
	s.clause.Set(clause.VALUES, recordsValues...)
	sqls, sqlArgs := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err = s.Raw(sqls, sqlArgs...).Exec()

	return
}
