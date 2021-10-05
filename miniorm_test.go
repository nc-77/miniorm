package miniorm

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dsn = "root:123456@tcp(localhost)/test"
	db  *DB
)

type Student struct {
	Id   uint `miniorm:"PRIMARY KEY"`
	Name string
	Age  int
}

func TestOpen(t *testing.T) {
	if _, err := Open("mysql", dsn); err != nil {
		t.Fatal()
	}
}

func TestMain(m *testing.M) {
	db, _ = Open("mysql", dsn)
	// run tests
	code := m.Run()
	os.Exit(code)
}

func TestDB_Table(t *testing.T) {
	if err := db.Model(&Student{}).CreateTable().Error(); err != nil {
		t.Fatal(err)
	}
	if exist := db.Model(&Student{}).HasTable(); !exist {
		t.Fatal("hasTable failed,excepted true,got false")
	}
	if err := db.Model(&Student{}).DropTable().Error(); err != nil {
		t.Fatal(err)
	}
}
