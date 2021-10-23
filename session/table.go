package session

import (
	"fmt"
	"strings"
)

func (s *Session) CreateTable() (err error) {
	defer func() {
		s.tableLast(0, err)
	}()

	table := s.RefTable()
	fields := make([]string, 0)
	for _, field := range table.Fields {
		fields = append(fields, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	args := strings.Join(fields, ",")
	_, err = s.Raw(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (%s)", table.Name, args)).Exec()

	return
}

func (s *Session) DropTable() (err error) {
	defer func() {
		s.tableLast(0, err)
	}()

	_, err = s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", s.RefTable().Name)).Exec()

	return
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExistSql(s.RefTable().Name)

	var result string
	row := s.Raw(sql, args...).QueryRow()
	_ = row.Scan(&result)

	return result == s.RefTable().Name
}
