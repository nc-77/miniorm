# miniorm[![Go Reference](https://pkg.go.dev/badge/github.com/nc-77/miniorm.svg)](https://pkg.go.dev/github.com/nc-77/miniorm)![example workflow](https://github.com/nc-77/miniorm/actions/workflows/go.yml/badge.svg)

[English](https://github.com/nc-77/miniorm/blob/main/README.md) | 中文

miniorm 是一个轻量级，快速的orm框架。框架提供接口主要参考了GROM 以及 XORM ，旨在为用户提供一个易上手，更自由的框架。

该项目主要是为了Go语言熟悉以及进阶，**请勿用于生产环境**。

## 特征

- 简洁，直观的API提供
- 更加自由，用户可控度更高的SQL接口提供
- 表的创建，删除以及主键的设置
- 记录的增删更查
- 支持钩子
- 支持事务
- 支持链式调用
- ...

## 使用

在项目中导入miniorm模块

```go
go get -u github.com/nc-77/miniorm
```

## 快速开始

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

更多示范请查看 [example.go](https://github.com/nc-77/miniorm/blob/main/_example/example.go) 

## 代办事项

- 更多的API提供
- 更多的数据库支持
- 完善相关文档以及注释
- ...

## 感谢

- [GeeORM教程](https://geektutu.com/post/geeorm.html)

