package schema

import (
	"miniorm/dialect"
	"testing"
)

type (
	student struct {
		Id uint
		Person
		Addr Address
	}
	Person struct {
		Name string
		Age  int
	}
	Address struct {
		City string
	}
)

type Student struct {
	Id   uint `miniorm:"PRIMARY KEY"`
	Name string
	Age  int
	any  string
}

var mysql, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	t.Run("normal struct", func(t *testing.T) {
		schema := Parse(&Student{}, mysql)
		if schema.Name != "Student" || len(schema.Fields) != 3 {
			t.Fatal("failed to parse Student struct")
		}
		field, ok := schema.GetField("Id")
		if !ok {
			t.Fatal("failed to parse id field")
		} else if field.Tag != "PRIMARY KEY" {
			t.Fatal("failed to parse PRIMARY KEY tag")
		} else if field.Type != "integer" {
			t.Fatal("failed to parse id field type")
		}
	})

	t.Run("embedded struct", func(t *testing.T) {
		schema := Parse(student{}, mysql)
		if schema.Name != "student" || len(schema.Fields) != 4 {
			t.Fatal("failed to parse student struct")
		}
	})
}
