package csv

import (
	"fmt"
)

type CsvError struct {
	message string
	field   string
}

func (e CsvError) Error() string {
	if e.field != "" {
		return fmt.Sprintf(e.message, e.field)
	}

	return e.message
}

func (e CsvError) SetField(fieldName string) {
	e.field = fieldName
}

var UnknownField = CsvError{message: "unknown field: %s"}
