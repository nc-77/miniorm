package session

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	s *Session
)

func newSession() *Session {
	dsn := "root:123456@tcp(localhost)/test"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return NewSession(db)
}

func TestMain(m *testing.M) {
	s = newSession()
	// create student table ans insert some rows for testing
	_, _ = s.db.Exec("CREATE  TABLE `student`(id int,name text,age text)")
	_, _ = s.db.Query("Insert INTO `student` VALUES(1,'nic',20)")
	_, _ = s.db.Query("INSERT INTO `student` VALUES(2,'zhangSan',18)")
	_, _ = s.db.Query("INSERT INTO `student` VALUES(3,'LiSi',19)")

	// run tests
	code := m.Run()
	// drop student table
	_, _ = s.db.Exec("DROP TABLE `student`")
	os.Exit(code)
}

func TestSession_Exec(t *testing.T) {

	if _, err := s.Raw("UPDATE `student` SET age = ? WHERE name = ?", 20, "nic").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("DELETE FROM `student` WHERE name = ?", "nic").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("INSERT INTO `student` VALUES(?,?,?)", 100, "nic", 24).Exec(); err != nil {
		t.Fatal()
	}
}

func TestSession_Query(t *testing.T) {

	if _, err := s.Raw("SELECT id FROM `student` WHERE name = ?", "nic").Query(); err != nil {
		t.Fatal()
	}
}

func TestSession_QueryRow(t *testing.T) {
	if _, err := s.Raw("SELECT id FROM `student` WHERE name = ?", "nic").QueryRow(); err != nil {
		t.Fatal()
	}

}
