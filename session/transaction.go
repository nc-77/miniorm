package session

import (
	"database/sql"
	"github.com/nc-77/miniorm/log"
)

type commDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var _ commDB = (*sql.DB)(nil)
var _ commDB = (*sql.Tx)(nil)

func (s *Session) DB() commDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	var tx *sql.Tx
	if tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	s.tx = tx
	return
}

func (s *Session) Commit() (err error) {
	defer func() {
		s.tx = nil
	}()
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *Session) Rollback() (err error) {
	defer func() {
		s.tx = nil
	}()
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	return
}
