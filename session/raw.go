package session

import (
	"database/sql"
	"fmt"
	"miniorm/dialect"
	"miniorm/log"
	"miniorm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlArgs  []interface{}
	dialect  dialect.Dialect // set when init
	refTable *schema.Schema  // set when model called
	err      error           // set when exec,query,queryRow go wrong
}

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		sqlArgs: make([]interface{}, 0),
		dialect: dialect,
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
	s.err = nil
	s.sql.Reset()
	s.sqlArgs = make([]interface{}, 0)
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sqlArgs = append(s.sqlArgs, args...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.clear()
	log.Info(s.toSql())
	if result, err = s.db.Exec(s.sql.String(), s.sqlArgs...); err != nil {
		s.err = err
		log.Error(err)
	}
	return
}

func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.clear()
	log.Info(s.toSql())
	if rows, err = s.db.Query(s.sql.String(), s.sqlArgs...); err != nil {
		s.err = err
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() (row *sql.Row) {
	defer s.clear()
	log.Info(s.toSql())
	row = s.db.QueryRow(s.sql.String(), s.sqlArgs...)
	s.err = row.Err()
	return
}

// Model set session's refTable
func (s *Session) Model(v interface{}) *Session {
	// 当前refTable为空或与指定新refTable时更新
	if s.refTable == nil || reflect.TypeOf(v) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(v, s.dialect)
	}
	return s
}

// RefTable return session's refTable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

func (s *Session) Err() error {
	return s.err
}
