package main

import (
	"fmt"
	"github.com/nc-77/miniorm"

	_ "github.com/go-sql-driver/mysql" // important!
)

type User struct {
	Id       uint `miniorm:"PRIMARY KEY"`
	Name     string
	Age      int
	Password string
	IsValid  bool
}

// AfterFirst hook.It will be called after calling first().
// Miniorm support hooks including before/after First/Insert/Delete. You can use them by adding the methods shown below.
func (u *User) AfterFirst() error {
	u.Password = "******"
	return nil
}

func (u *User) BeforeInsert() error {
	u.Id += 1000
	return nil
}

func main() {

	db, err := miniorm.Open("mysql", "root:123456@tcp(localhost)/test")
	if err != nil {
		panic(err)
	}
	// Delete Table
	db.DropTable(User{})
	// Create Table
	db.CreateTable(User{})
	// CheckTable
	if ok := db.HasTable(User{}); !ok {
		return
	}

	// Records
	users := []User{{Id: 1, Name: "nic", Age: 11, IsValid: true}, {Id: 2, Name: "nc-77", Age: 22, IsValid: true}, {Id: 3, Name: "nc-77", Age: 11, IsValid: false}}
	// Create
	db.Create(&users[0], &users[1], &users[2])
	db.CreateMany(&users)

	var foundUsers []User
	var foundUser User
	// Find all records
	result := db.Find(&foundUsers)
	fmt.Println(result.RowsAffected)

	// Find specific records
	db.Where("IsValid = ?", false).Limit(1).Find(&foundUsers)
	db.Where("IsValid = ? and Age = ?", true, 22).First(&foundUser)

	// Update
	users[1].Age++
	users[1].IsValid = false
	// default save all fields even if it is zero
	db.Update(users[1], users[2])
	// update specific fields
	db.Model(users[2]).Update(map[string]interface{}{"Age": users[2].Age + 1, "IsValid": false})
	// add where condition
	db.Model(User{}).Where("id = ?", 3).Update(map[string]interface{}{"Age": users[2].Age + 1, "IsValid": false})

	// Delete
	db.Model(User{}).Where("id = ?", 3).Delete()

	// Raw SQL and SQL builder
	db.Raw("SELECT Id,Name,Age FROM `User` WHERE IsValid = ?", true).Query()
	db.Raw("INSERT INTO `User` (Id,Name,Age,IsValid) VALUES (1010,'nic',12,true)").Exec()

	// Transactions
	db.Transaction(func(db *miniorm.DB) (err error) {
		users[0].Age++
		err = db.Update(users[0]).Error
		err = db.Create(users[0]).Error
		// if err is not nil, it will rollback and return err, otherwise commit
		return
	})

}
