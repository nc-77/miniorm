package schema

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/nc-77/miniorm/dialect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	PrimaryKey []string
	Fields     []*Field
	FieldsName []string
	FieldsMap  map[string]*Field
}

// Parse parses any struct to schema
func Parse(v interface{}, d dialect.Dialect) (s *Schema) {
	model := reflect.Indirect(reflect.ValueOf(v)).Type()

	s = &Schema{
		Model:      v,
		Name:       model.Name(),
		PrimaryKey: make([]string, 0),
		Fields:     make([]*Field, 0),
		FieldsName: make([]string, 0),
		FieldsMap:  make(map[string]*Field),
	}

	var parseField func(reflect.StructField)
	parseField = func(p reflect.StructField) {
		if ast.IsExported(p.Name) {
			if p.Type.Kind() != reflect.Struct {
				field := &Field{
					Name: p.Name,
					Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				}
				if tag, ok := p.Tag.Lookup("miniorm"); ok {
					field.Tag = tag
					if tag == "PRIMARY KEY" {
						s.PrimaryKey = append(s.PrimaryKey, field.Name)
					}
				}
				if _, ok := s.FieldsMap[field.Name]; ok {
					panic(fmt.Sprintf("%s field is repeated in the same structure.\n", field.Name))
				}
				s.Fields = append(s.Fields, field)
				s.FieldsName = append(s.FieldsName, field.Name)
				s.FieldsMap[field.Name] = field
			} else {
				model := p.Type
				for i := 0; i < model.NumField(); i++ {
					parseField(model.Field(i))
				}
			}
		}
	}
	for i := 0; i < model.NumField(); i++ {
		p := model.Field(i)
		parseField(p)
	}

	return
}

func (s *Schema) GetField(name string) (field *Field, ok bool) {
	field, ok = s.FieldsMap[name]
	return
}
