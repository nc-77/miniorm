package miniorm

import (
	"database/sql"
	"errors"
	"fmt"
	"miniorm/dialect"

	"miniorm/log"
	"miniorm/session"
)

type DB struct {
	db      *sql.DB
	session *session.Session
}

// Open initialize db session based on dsn
func Open(drive, dsn string) (*DB, error) {
	db, err := sql.Open(drive, dsn)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	d, ok := dialect.GetDialect(drive)
	if !ok {
		err := errors.New(fmt.Sprintf("dialect %s Not Found", drive))
		log.Error(err)
		return nil, err
	}
	retDB := &DB{
		db:      db,
		session: session.NewSession(db, d),
	}
	log.Infof("Connect %s success", drive)
	return retDB, nil
}

// Error return db error
func (db *DB) Error() error {
	return db.session.Err()
}

// Model specify the model you would like to run db operations
func (db *DB) Model(value interface{}) *DB {
	_ = db.session.Model(value)
	return db
}

// HasTable create db table based on Model
func (db *DB) HasTable() bool {
	return db.session.HasTable()
}

// CreateTable create db table based on Model
func (db *DB) CreateTable() *DB {
	_ = db.session.CreateTable()
	return db
}

// DropTable drop db table based on Model
func (db *DB) DropTable() *DB {
	_ = db.session.DropTable()
	return db
}
