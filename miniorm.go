package miniorm

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/nc-77/miniorm/dialect"
	"github.com/nc-77/miniorm/log"
	"github.com/nc-77/miniorm/session"
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

// HasTable check db table exist based on Value struct
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
	destType := reflect.Indirect(reflect.ValueOf(values)).Type().Elem()
	model := reflect.New(destType).Interface()
	_ = db.session.Model(model).FindRecords(values)
	return db.session.Result()
}

// First find the first record that match given conditions
func (db *DB) First(value interface{}) *session.Result {
	_ = db.session.Model(value).FirstRecord(value)
	return db.session.Result()
}

// Create insert records based on values
func (db *DB) Create(values ...interface{}) *session.Result {
	if len(values) == 0 {
		panic("should pass at least one argument to db.Create")
	}
	_ = db.session.Model(values).CreateRecords(values...)
	return db.session.Result()
}

// CreateMany insert records based on values
func (db *DB) CreateMany(values interface{}) *session.Result {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	dest := make([]interface{}, destSlice.Len())
	for i := range dest {
		dest[i] = destSlice.Index(i).Addr().Interface()
	}
	_ = db.session.Model(destSlice.Index(0).Interface()).CreateRecords(dest...)
	return db.session.Result()
}

// Delete delete records after setting where
func (db *DB) Delete() *session.Result {
	_ = db.session.DeleteRecords()
	return db.session.Result()
}

// Update update all fields including zero when values is struct or []struct.
// Update update specified fields when values is map[string]interface{}.
func (db *DB) Update(values ...interface{}) *session.Result {
	if len(values) == 0 {
		panic("should pass at least one argument to db.Update")
	}
	v := reflect.Indirect(reflect.ValueOf(values[0]))
	switch v.Kind() {
	case reflect.Struct:
		{
			for _, value := range values {
				_ = db.session.Model(value).UpdateRecords(value)
			}
		}
	case reflect.Map:
		{
			_ = db.session.UpdateRecords(values[0])
		}
	}
	return db.session.Result()
}

// Where add conditions
func (db *DB) Where(desc string, values ...interface{}) *DB {
	var vars []interface{}
	db.session.Where(append(append(vars, desc), values...)...)
	return db
}

// Limit add conditions
func (db *DB) Limit(limit int) *DB {
	db.session.Limit(limit)
	return db
}

// OrderBy add conditions
func (db *DB) OrderBy(values ...interface{}) *DB {
	db.session.OrderBy(values...)
	return db
}

func (db *DB) Raw(sql string, args ...interface{}) *DB {
	_ = db.session.Raw(sql, args...)
	return db
}
func (db *DB) Exec() (result sql.Result, err error) {
	return db.session.Exec()
}

func (db *DB) Query() (rows *sql.Rows, err error) {
	return db.session.Query()
}

func (db *DB) QueryRaw() (row *sql.Row) {
	return db.session.QueryRow()
}
