package sqlstore

import "fmt"

type Dialect string

const (
	DialectMySQL Dialect = "mysql"
)

func ParseDialect(value string) (Dialect, error) {
	if value == "" || value == string(DialectMySQL) {
		return DialectMySQL, nil
	}
	return "", fmt.Errorf("unsupported dialect %q", value)
}
