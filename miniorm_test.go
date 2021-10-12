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

type User struct {
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
	if result := db.CreateTable(User{}); result.Error != nil {
		t.Fatal(result.Error)
	}
	if exist := db.HasTable(&User{}); !exist {
		t.Fatal("hasTable failed,excepted true,got false")
	}
	if result := db.DropTable(&User{}); result.Error != nil {
		t.Fatal(result.Error)
	}
}

func TestDB_Find(t *testing.T) {

}
