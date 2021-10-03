package miniorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestOpenDB(t *testing.T) {
	dsn := "root:123456@tcp(localhost)/test"
	if _, err := OpenDB("mysql", dsn); err != nil {
		t.Fatal()
	}
}
