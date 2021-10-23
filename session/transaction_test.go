package session

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSession_transaction(t *testing.T) {
	var stu Student
	stu1 := Student{Id: 10, Name: "nc-77", Age: 21}
	stu2 := Student{Id: 11, Name: "nc", Age: 22}
	if err := s.Begin(); err != nil {
		t.Fatal(err)
	}
	if err := s.Model(Student{}).CreateRecords(stu1, stu2); err != nil {
		t.Fatal(err)
	}
	if err := s.Rollback(); err != nil {
		t.Fatal(err)
	}
	if err := s.Model(Student{}).Where("id = ?", 10).FirstRecord(&stu); err != nil {
		t.Fatal(err)
	}
	if s.result.Error != ErrRecordNotFound {
		t.Fatal("rollback failed")
	}
}
