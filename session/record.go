package session

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/nc-77/miniorm/clause"
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
	sqls, sqlArgs := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

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
	s.CallMethod(BeforeFirst)

	s.clause.Set(clause.LIMIT, 1)

	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err = s.FindRecords(destSlice.Addr().Interface()); err != nil {
		return
	}
	if destSlice.Len() > 0 {
		dest.Set(destSlice.Index(0))
	}

	s.CallMethod(AfterFirst)
	return
}

// CreateRecord insert records into database by parsing value
func (s *Session) CreateRecords(values ...interface{}) (err error) {
	var result sql.Result
	defer func() {
		var affected int64
		if result != nil {
			affected, _ = result.RowsAffected()
		}
		s.recordLast(affected, err)
	}()
	// beforeInsert hook
	s.CallMethod(BeforeInsert)

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

	// afterInsert hook
	s.CallMethod(AfterInsert)

	return
}

// DeleteRecords delete records from database by setting where
func (s *Session) DeleteRecords() (err error) {
	var result sql.Result
	defer func() {
		var affected int64
		if result != nil {
			affected, _ = result.RowsAffected()
		}
		s.recordLast(affected, err)
	}()
	// beforeDelete hook
	s.CallMethod(BeforeDelete)

	table := s.RefTable()
	s.clause.Set(clause.DELETE, table.Name)
	sqls, sqlArgs := s.clause.Build(clause.DELETE, clause.WHERE)

	result, err = s.Raw(sqls, sqlArgs...).Exec()

	// afterDelete hook
	s.CallMethod(AfterDelete)

	return
}

// UpdateRecords update all fields even if the field is zero
func (s *Session) UpdateRecords(values ...interface{}) (err error) {
	var result sql.Result
	defer func() {
		var affected int64
		if result != nil {
			affected, _ = result.RowsAffected()
		}
		s.recordLast(affected, err)
	}()
	// beforeUpdate hook
	//s.CallMethod(BeforeUpdate)

	table := s.RefTable()
	pks := make([]string, len(table.PrimaryKey))
	for i := range pks {
		pks[i] = table.PrimaryKey[i] + " = ?"
	}

	for _, value := range values {
		pkValues := make([]interface{}, len(table.PrimaryKey))
		v := reflect.Indirect(reflect.ValueOf(value))
		switch v.Kind() {
		case reflect.Struct: // 传入结构体时根据结构体中主键设置where 主键
			{
				for i := range pkValues {
					pk := table.PrimaryKey[i]
					pkValues[i] = v.FieldByName(pk).Interface()
				}
				s.clause.Set(clause.UPDATE, table.Name, value)
			}
		case reflect.Map: // 传入map时根据之前的model设置where 主键
			{
				for i := range pkValues {
					pk := table.PrimaryKey[i]
					pkValues[i] = reflect.Indirect(reflect.ValueOf(table.Model)).FieldByName(pk).Interface()
				}
				s.clause.Set(clause.UPDATE, table.Name, value)
			}
		}
		// 用户未指定where时默认设置主键
		if !s.clause.Exist(clause.WHERE) {
			s.clause.Set(clause.WHERE, append(append([]interface{}{}, strings.Join(pks, "and")), pkValues...)...)
		}
		sqls, sqlArgs := s.clause.Build(clause.UPDATE, clause.WHERE)
		result, err = s.Raw(sqls, sqlArgs...).Exec()
	}
	// afterUpdate hook
	//s.CallMethod(AfterUpdate)

	return
}

func (s *Session) Where(values ...interface{}) *Session {
	s.clause.Set(clause.WHERE, values...)
	return s
}

func (s *Session) Limit(values ...interface{}) *Session {
	s.clause.Set(clause.LIMIT, values...)
	return s
}

func (s *Session) OrderBy(values ...interface{}) *Session {
	s.clause.Set(clause.ORDERBY, values...)
	return s
}
