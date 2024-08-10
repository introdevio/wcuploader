package uploader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/introdevio/wcuploader/internal/product"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ProductLoader struct {
	csvPath     string
	rootDir     string
	tagFieldMap map[string]string
}

func NewProductLoader(csvPath, rootDir string) *ProductLoader {
	return &ProductLoader{
		csvPath: csvPath,
		rootDir: rootDir,
	}
}

func (pl *ProductLoader) LoadFromCsv(hasHeader bool) ([]product.Product, []error) {

	file, err := os.Open(pl.csvPath)

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			log.Println("Failed to close file", file.Name())
		}
	}(file)

	if err != nil {
		return nil, []error{err}
	}

	reader := csv.NewReader(file)
	headers, err := loadHeaders(reader, hasHeader)

	if err != nil {
		return nil, []error{err}
	}

	var errs []error
	var items []product.Product

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

		st := product.Product{}
		tagFieldMap := loadFieldTagMap(st)
		for i, value := range line {
			reflected := reflect.ValueOf(&st).Elem()
			tagName := strings.ToLower(headers[i])
			if _, exists := tagFieldMap[tagName]; !exists {
				errs = append(errs, errors.New("unknown field: "+headers[i]))
				continue
			}
			field := reflected.FieldByName(tagFieldMap[tagName])
			if value == "" {
				continue
			}
			err := setField(field, value)
			if err != nil {
				errs = append(errs, err)
			}
		}
		items = append(items, st)
	}
	return items, errs
}

func loadHeaders(reader *csv.Reader, hasHeader bool) ([]string, error) {
	var headers []string
	var err error
	if hasHeader {
		headers, err = reader.Read()
		if err != nil {
			return nil, err
		}
	}
	return headers, nil
}

func loadFieldTagMap(input any) map[string]string {
	tagFieldMap := make(map[string]string)

	t := reflect.TypeOf(input)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		tagFieldMap[tag] = field.Name
	}

	return tagFieldMap
}

func setField(field reflect.Value, value string) error {
	switch field.Type().String() {
	case "string":
		field.SetString(value)
	case "int":
		v, err := strconv.Atoi(value)
		if err != nil {
			msg := fmt.Sprintf("Error converting %s to int", value)
			return errors.New(msg)
		}
		field.SetInt(int64(v))
	case "float32":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			msg := fmt.Sprintf("Error converting %s to float", value)
			return errors.New(msg)
		}
		field.SetFloat(v)
	}
	return nil
}
