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

func TestSession_DeleteRecords(t *testing.T) {
	if err := s.Model(Student{}).Where("Name = ?", "nic").DeleteRecords(); err != nil {
		t.Fatal(err)
	}
	if s.Result().RowsAffected != 1 {
		t.Fatalf("rowsAffected failed,excepted %v,got %v", 1, s.Result().RowsAffected)
	}

}

func TestSession_FirstRecord(t *testing.T) {
	var stu Student
	if err := s.Model(Student{}).FirstRecord(&stu); err != nil {
		t.Fatal(err)
	}
	if s.Result().RowsAffected != 1 {
		t.Fatalf("rowsAffected failed,excepted %v,got %v", 1, s.Result().RowsAffected)
	}

}

func TestSession_UpdateRecords(t *testing.T) {
	stu1 := Student{Id: 1, Name: "nic", Age: 19}
	stu2 := Student{Id: 2, Name: "nc", Age: 18}

	if err := s.Model(Student{}).UpdateRecords(&stu1, &stu2); err != nil {
		t.Fatal(err)
	}
	if err := s.Model(&stu1).UpdateRecords(map[string]interface{}{"Age": stu1.Age + 1}); err != nil {
		t.Fatal(err)
	}
	t.Log(stu1)
}
