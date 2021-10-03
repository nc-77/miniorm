package session

import (
	"database/sql"
	"fmt"
	"miniorm/log"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlArgs []interface{}
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		db:      db,
		sqlArgs: make([]interface{}, 0),
	}
}

// toSql combine sql and sqlArgs into a complete SQL, only used in log
func (s *Session) toSql() (sql string) {
	sql = s.sql.String()
	for i := range s.sqlArgs {
		sql = strings.Replace(sql, "?", fmt.Sprint(s.sqlArgs[i]), 1)
	}
	return
}

func (s *Session) clear() {
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
		log.Error(err)
	}
	return
}

func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.clear()
	log.Info(s.toSql())
	if rows, err = s.db.Query(s.sql.String(), s.sqlArgs...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() (row *sql.Row, err error) {
	defer s.clear()
	log.Info(s.toSql())
	row = s.db.QueryRow(s.sql.String(), s.sqlArgs...)
	if err = row.Err(); err != nil {
		log.Error(err)
	}
	return
}
