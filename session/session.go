package session

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/nc-77/miniorm/clause"
	"github.com/nc-77/miniorm/dialect"
	"github.com/nc-77/miniorm/log"
	"github.com/nc-77/miniorm/schema"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlArgs  []interface{}
	dialect  dialect.Dialect // set when init
	refTable *schema.Schema  // set when model is called
	clause   *clause.Clause  // combines the clauses into a complete sql and sqlArgs
	result   *Result         // set when exec,query,queryRow are called
}

type Result struct {
	Error        error
	RowsAffected int64
}

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		sqlArgs: make([]interface{}, 0),
		dialect: dialect,
		clause:  new(clause.Clause),
		result:  new(Result),
	}
}

// toSql combine sql and sqlArgs into a complete SQL, only used in log
func (s *Session) toSql() (sql string) {
	sql = s.sql.String()
	for i := range s.sqlArgs {
		if reflect.TypeOf(s.sqlArgs[i]).Kind() == reflect.String {
			sql = strings.Replace(sql, "?", fmt.Sprintf("'%v'", s.sqlArgs[i]), 1)
		} else {
			sql = strings.Replace(sql, "?", fmt.Sprint(s.sqlArgs[i]), 1)
		}
	}
	return
}

func (s *Session) clear() {
	s.clause.Clear()
	s.sql.Reset()
	s.sqlArgs = make([]interface{}, 0)
}

func (s *Session) Result() *Result {
	return s.result
}

// last set session's result after crud records
func (s *Session) recordLast(affected int64, err error) {
	s.result = &Result{
		Error:        err,
		RowsAffected: affected,
	}
	if err == nil && affected == 0 {
		s.result.Error = ErrRecordNotFound
	}
}

// last set session's result after crud table
func (s *Session) tableLast(affected int64, err error) {
	s.result = &Result{
		Error:        err,
		RowsAffected: affected,
	}
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sqlArgs = append(s.sqlArgs, args...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.clear()
	log.Info(s.toSql())
	result, err = s.db.Exec(s.sql.String(), s.sqlArgs...)
	if err != nil {
		log.Error(err)
	}

	return
}

func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.clear()
	log.Info(s.toSql())
	rows, err = s.db.Query(s.sql.String(), s.sqlArgs...)
	if err != nil {
		log.Error(err)
	}

	return
}

func (s *Session) QueryRow() (row *sql.Row) {
	defer s.clear()
	log.Info(s.toSql())
	row = s.db.QueryRow(s.sql.String(), s.sqlArgs...)

	return
}

// Model set session's refTable
func (s *Session) Model(v interface{}) *Session {
	// 当前refTable为空或与指定新refTable时更新
	if s.refTable == nil || reflect.TypeOf(v) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(v, s.dialect)
	}
	s.refTable.Model = v
	return s
}

// RefTable return session's refTable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}
