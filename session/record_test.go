package session

import (
	"testing"
)

func TestSession_FindRecords(t *testing.T) {
	var stu []Student
	err := s.Model(Student{}).FindRecords(&stu)
	if err != nil {
		t.Fatal(err)
	}
	if s.result.RowsAffected != 3 {
		t.Fatalf("rowsAffected exceptd %v,got %v", 3, s.result.RowsAffected)
	}
	if len(stu) != 3 {
		t.Fatalf("stu len failed exceptd %v,got %v", 3, len(stu))
	}
	t.Log(stu)
}

func TestSession_CreateRecord(t *testing.T) {

	stu1 := Student{Id: 10, Name: "nc-77", Age: 21}
	stu2 := Student{Id: 11, Name: "nc", Age: 22}
	if err := s.Model(stu1).CreateRecords(stu1, stu2); err != nil {
		t.Fatal(err)
	}
	if s.Result().RowsAffected != 2 {
		t.Fatalf("rowsAffected failed,excepted %v,got %v", 2, s.Result().RowsAffected)
	}
	var stus []Student
	if err := s.Model(Student{}).FindRecords(&stus); err != nil {
		t.Fatal(err)
	}
	if s.Result().RowsAffected != 5 {
		t.Fatalf("find all records failed,excepted %v,got %v", 5, s.Result().RowsAffected)
	}
}
