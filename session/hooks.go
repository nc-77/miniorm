package session

import (
	"github.com/nc-77/miniorm/log"
	"reflect"
)

const (
	BeforeFirst  = "BeforeFirst"
	AfterFirst   = "AfterFirst"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"

	//BeforeUpdate = "BeforeUpdate"
	//AfterUpdate  = "AfterUpdate"
)

// CallMethod is in order to implement the hook
func (s *Session) CallMethod(method string) {
	v := reflect.ValueOf(s.RefTable().Model)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			m := reflect.ValueOf(v.Index(i).Interface()).MethodByName(method)
			if m.IsValid() {
				if outs := m.Call([]reflect.Value{}); len(outs) > 0 {
					if err, ok := outs[0].Interface().(error); ok {
						log.Error(err)
					}
				}
			}
		}
	default:
		m := v.MethodByName(method)
		if m.IsValid() {
			if outs := m.Call([]reflect.Value{}); len(outs) > 0 {
				if err, ok := outs[0].Interface().(error); ok {
					log.Error(err)
				}
			}
		}
	}

}
