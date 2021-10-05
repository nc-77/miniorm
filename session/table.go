package session

import (
	"fmt"
	"strings"
)

func (s *Session) CreateTable() error {
	table := s.RefTable()
	fields := make([]string, 0)
	for _, field := range table.Fields {
		fields = append(fields, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	args := strings.Join(fields, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE `%s` (%s)", table.Name, args)).Exec()

	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExistSql(s.RefTable().Name)

	var result string
	row := s.Raw(sql, args...).QueryRow()
	_ = row.Scan(&result)

	return result == s.RefTable().Name
}
