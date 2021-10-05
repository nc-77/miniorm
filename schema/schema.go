package schema

import (
	"fmt"
	"go/ast"
	"miniorm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model    interface{}
	Name     string
	Fields   []*Field
	fieldMap map[string]*Field
}

// Parse parses any struct to schema
func Parse(v interface{}, d dialect.Dialect) (s *Schema) {
	model := reflect.Indirect(reflect.ValueOf(v)).Type()

	s = &Schema{
		Model:    v,
		Name:     model.Name(),
		Fields:   make([]*Field, 0),
		fieldMap: make(map[string]*Field),
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
				}
				if _, ok := s.fieldMap[field.Name]; ok {
					panic(fmt.Sprintf("%s field is repeated in the same structure.\n", field.Name))
				}
				s.Fields = append(s.Fields, field)
				s.fieldMap[field.Name] = field
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
	field, ok = s.fieldMap[name]
	return
}
