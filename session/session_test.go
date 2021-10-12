package session

import (
	"database/sql"
	"miniorm/dialect"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	s *Session
)

var mysql, _ = dialect.GetDialect("mysql")

func newSession() *Session {
	dsn := "root:123456@tcp(localhost)/test"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return NewSession(db, mysql)
}

func TestMain(m *testing.M) {
	s = newSession()
	// create student table ans insert some rows for testing
	_, _ = s.db.Exec("CREATE  TABLE Student(Id int,Name text,Age text)")
	_, _ = s.db.Query("INSERT INTO `Student` VALUES(1,'nic',18)")
	_, _ = s.db.Query("INSERT INTO `Student` VALUES(2,'zhangSan',18)")
	_, _ = s.db.Query("INSERT INTO `Student` VALUES(3,'LiSi',19)")

	// run tests
	code := m.Run()
	// drop student table
	_, _ = s.db.Exec("DROP TABLE IF EXISTS  `Student`")
	os.Exit(code)
}

func TestSession_Exec(t *testing.T) {

	if _, err := s.Raw("UPDATE `Student` SET Age = ? WHERE Name = ?", 20, "nic").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("DELETE FROM `Student` WHERE Name = ?", "nic").Exec(); err != nil {
		t.Fatal()
	}
	if _, err := s.Raw("INSERT INTO `Student` VALUES(?,?,?)", 100, "nic", 24).Exec(); err != nil {
		t.Fatal()
	}
}

func TestSession_Query(t *testing.T) {

	if _, err := s.Raw("SELECT Id FROM `Student` WHERE Name = ?", "nic").Query(); err != nil {
		t.Fatal()
	}
}

func TestSession_QueryRow(t *testing.T) {
	if row := s.Raw("SELECT Id FROM `Student` WHERE Name = ?", "nic").QueryRow(); row.Err() != nil {
		t.Fatal()
	}

}
