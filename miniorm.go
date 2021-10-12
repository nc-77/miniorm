package miniorm

import (
	"database/sql"
	"fmt"
	"miniorm/dialect"
	"reflect"

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
		err := fmt.Errorf("dialect %s Not Found", drive)
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

// Model specify the model you would like to run db operations
func (db *DB) Model(value interface{}) *DB {
	_ = db.session.Model(value)
	return db
}

// HasTable create db table based on Value struct
func (db *DB) HasTable(value interface{}) bool {
	return db.session.Model(value).HasTable()
}

// CreateTable create db table based on Value struct
func (db *DB) CreateTable(value interface{}) *session.Result {
	_ = db.session.Model(value).CreateTable()
	return db.session.Result()
}

// DropTable drop db table based on Value struct
func (db *DB) DropTable(value interface{}) *session.Result {
	_ = db.session.Model(value).DropTable()
	return db.session.Result()
}

// Find find records based on values
func (db *DB) Find(values interface{}) *session.Result {
	model := reflect.ValueOf(values).Elem()
	_ = db.session.Model(model).FindRecords(values)
	return db.session.Result()
}
