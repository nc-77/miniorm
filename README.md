# miniorm[![Go Reference](https://pkg.go.dev/badge/github.com/nc-77/miniorm.svg)](https://pkg.go.dev/github.com/nc-77/miniorm)![example workflow](https://github.com/nc-77/miniorm/actions/workflows/go.yml/badge.svg)

English | [中文](https://github.com/nc-77/miniorm/blob/main/README-zh.md)

Miniorm is a lightful and fast orm framework.The API mainly refers to Gorm and Xorm, and the purpose is to provide users with a more easy and free framework.

This project began as a way to learn Golang. **Don't use in any production environment**.

## Features

- Easy API supported
- More freedom SQL API supported
- Table CRUD and primary key settings
- Record CRUD
- Transaction
- Hooks
- Support for chain calls
- ...

## Usage

Import miniorm in the applicaton

```go
go get -u github.com/nc-77/miniorm
```

## QuickStart

```go
package main
import (
	"fmt"
	"github.com/nc-77/miniorm"
	_ "github.com/go-sql-driver/mysql" 
)
type User struct {
	Id      uint `miniorm:"PRIMARY KEY"`
	Name    string
	Age     int
	IsValid bool
}
func main(){
    db, err := miniorm.Open("mysql", "root:123456@tcp(localhost)/test")
	if err != nil {
		panic(err)
	}
    // Create Table
	db.CreateTable(User{})
    
    // Records
	users := []User{{Id: 1, Name: "nic", Age: 11, IsValid: true}, {Id: 2, Name: "nc-77", Age: 22, IsValid: true}, {Id: 3, Name: "nc-77", Age: 11, IsValid: false}}
	// Create
	db.Create(users[0], users[1], users[2])
	db.CreateMany(users)
    
    // Find all records
	result := db.Find(&foundUsers)
	fmt.Println(result.RowsAffected)

	// Find specific records
	db.Where("IsValid = ?", false).Limit(1).Find(&foundUsers)
	db.Where("IsValid = ? and Age = ?", true, 22).First(&foundUser)
    
    ...
    
    return 
}
```

Please see [example.go](https://github.com/nc-77/miniorm/blob/main/_example/example.go) for more imformation

## Todo

- More API supported
- More Database supported
- More complete documentation and notes

## Thanks

- [GeeORM](https://geektutu.com/post/geeorm.html)
