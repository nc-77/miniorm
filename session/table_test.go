package session

import (
	"testing"
)

type Student struct {
	Id   uint `miniorm:"PRIMARY KEY"`
	Name string
	Age  int
}

func TestSession_Table(t *testing.T) {
	if exist := s.Model(&Student{}).HasTable(); !exist {
		t.Fatal("hasTable failed,excepted true,got false")
	}
	if err := s.Model(&Student{}).DropTable(); err != nil {
		t.Fatal(err)
	}
	if exist := s.Model(&Student{}).HasTable(); exist {
		t.Fatal("hasTable failed,excepted false,got true")
	}
	if err := s.Model(&Student{}).CreateTable(); err != nil {
		t.Fatal()
	}
}
