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
