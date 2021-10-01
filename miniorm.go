package miniorm

import (
	"database/sql"

	"miniorm/log"
	"miniorm/session"
)

type DB struct {
	db      *sql.DB
	session *session.Session
}

func OpenDB(drive, dsn string) (*DB, error) {
	db, err := sql.Open(drive, dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("Connect database success")
	retDB := &DB{
		db:      db,
		session: session.NewSession(db),
	}
	return retDB, nil
}
