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

func TestDB_Record(t *testing.T) {
	user0 := User{Id: 0, Name: "nic", Age: 18}
	user1 := &User{Id: 1, Name: "nc-77", Age: 20}
	// create table
	if result := db.CreateTable(User{}); result.Error != nil {
		t.Fatal(result.Error)
	}
	// insert records
	result := db.Create(user0, user1)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if result.RowsAffected != 2 {
		t.Fatalf("create records failed,rowsAffected excepted %v,got %v", 2, result.RowsAffected)
	}
	// find records
	var users []User
	result = db.Find(&users)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if users[0] != user0 || users[1] != *user1 {
		t.Log(users)
		t.Fatal("find records failed")
	}

}
