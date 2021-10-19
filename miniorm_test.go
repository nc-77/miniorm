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
	// find first record
	var firstUser User
	if result = db.OrderBy("Id", "DESC").First(&firstUser); result.Error != nil {
		t.Fatal(result.Error)
	}
	if firstUser != *user1 {
		t.Log(firstUser)
		t.Fatal("find first record failed")
	}
	// update records
	user0.Age++
	user1.Age++
	if result = db.Update(user0, user1); result.Error != nil {
		t.Fatal(result.Error)
	}
	if result = db.Model(user1).Update(map[string]interface{}{"Age": user1.Age + 1}); result.Error != nil {
		t.Fatal(result.Error)
	}
	if result = db.Model(&User{}).Where("Name = ?", "nc-77").Update(map[string]interface{}{"Age": 23}); result.Error != nil {
		t.Fatal(result.Error)
	}
	// delete records
	if result = db.Model(&User{}).Where("Name = ? and Age = ?", "nc-77", 23).Delete(); result.Error != nil {
		t.Fatal(result.Error)
	}
}
