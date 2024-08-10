package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"reflect"
)

func LoadFromCsv(filePath string, hasHeader bool) []error {

	file, err := os.Open(filePath)

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			log.Println("Failed to close file", file.Name())
		}
	}(file)

	if err != nil {
		return []error{err}
	}

	reader := csv.NewReader(file)
	var headers []string
	if hasHeader {
		headers, err = reader.Read()
		if err != nil {
			return []error{err}
		}
	}

	var errs []error

	for {
		line, err := reader.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("EOF")
				break
			}
			errs = append(errs, err)
			break
		}

		typ := reflect.TypeOf((*T)(nil)).Elem()
		st := reflect.New(typ).Elem().Interface().(*T)

		tagFieldMap := loadFieldTagMap(st)

		for i, value := range line {
			reflected := reflect.ValueOf(&st).Elem()
			if _, exists := tagFieldMap[headers[i]]; !exists {
				errs = append(errs, errors.New("unknown field: "+headers[i]))
				continue
			}
			field := reflected.FieldByName(tagFieldMap[headers[i]])
			err := setField(field, value)
			if err != nil {
				errs = append(errs, err)
			}
		}
		items = append(items, st)
	}
	return errs
}
