package session

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	s *Session
)

func newSession() *Session {
	dsn := "root:123456@tcp(localhost)/student"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return NewSession(db)
}

func TestMain(m *testing.M) {
	s = newSession()
	code := m.Run()
	os.Exit(code)
}

func TestSession_Exec(t *testing.T) {
	if _, err := s.Raw("CREATE TABLE IF NOT EXISTS User(id int,name text)").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("INSERT INTO User(id,name) VALUES (?,?)", 1, "nic").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("DROP TABLE User").Exec(); err != nil {
		t.Fatal()
	}
}
