package miniorm

import (
	"errors"
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

type UserWithHocks struct {
	User
	Password           string
	BfDelete, AtDelete bool
}

func (u *UserWithHocks) AfterFirst() error {
	u.Password = "******"
	return nil
}

func (u *UserWithHocks) BeforeInsert() error {
	u.Id += 1000
	return nil
}

func (u *UserWithHocks) BeforeDelete() error {
	u.BfDelete = true
	return nil
}

func (u *UserWithHocks) AfterDelete() error {
	u.AtDelete = true
	return nil
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
	defer func() {
		db.DropTable(&User{})
	}()
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

func TestDB_Hock(t *testing.T) {
	var u0, u1 UserWithHocks
	u1.Id = 1
	if err := db.CreateTable(UserWithHocks{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&u0, &u1).Error; err != nil {
		t.Fatal(err)
	}
	t.Run("beforeInsert && afterFind", func(t *testing.T) {
		var foundUser UserWithHocks
		if err := db.First(&foundUser).Error; err != nil {
			t.Fatal(err)
		}
		if foundUser.Id != 1000 || foundUser.Password != "******" {
			t.Fatal("first hooks failed")
		}
		var foundUsers []UserWithHocks
		if err := db.Find(&foundUsers).Error; err != nil {
			t.Fatal(err)
		}
		if foundUsers[0].Id != 1000 || foundUsers[1].Id != 1001 {
			t.Fatal("insert hooks failed")
		}
	})
	t.Run("before delete && after delete", func(t *testing.T) {
		var user UserWithHocks
		if err := db.Model(&user).Where("Id = ?", 1000).Delete().Error; err != nil {
			t.Fatal(err)
		}
		if !user.BfDelete || !user.AtDelete {
			t.Fatal("delete hooks failed")
		}
	})
	if err := db.DropTable(UserWithHocks{}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestDB_Transaction(t *testing.T) {
	defer func() {
		db.DropTable(&User{})
	}()
	user0 := User{Id: 0, Name: "nic", Age: 18}
	user1 := User{Id: 1, Name: "nc-77", Age: 20}
	db.CreateTable(&User{})

	t.Run("commit", func(t *testing.T) {
		if err := db.Transaction(func(db *DB) (err error) {
			if err = db.Create(&user0).Error; err != nil {
				return
			}
			if err = db.Create(&user1).Error; err != nil {
				return
			}
			return
		}); err != nil {
			t.Fatal(err)
		}
		var user []User
		db.Find(&user)
		if user[0] != user0 || user[1] != user1 {
			t.Fatal("commit failed")
		}
	})

	t.Run("rollback", func(t *testing.T) {
		if err := db.Transaction(func(db *DB) error {
			db.Model(user0).Update(map[string]interface{}{"Id": user0.Id + 1000})
			db.Model(user1).Update(map[string]interface{}{"Age": user1.Age + 1})
			return errors.New("error")
		}); err == nil {
			t.Fatal("rollback failed")
		}
		var user []User
		db.Find(&user)
		if user[0] != user0 || user[1] != user1 {
			t.Fatal("rollback failed")
		}
	})

}
