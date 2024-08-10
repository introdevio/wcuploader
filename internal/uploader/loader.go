package uploader

import (
	"encoding/csv"
	"errors"
	"github.com/introdevio/wcuploader/internal/product"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
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

func (pl *ProductLoader) LoadFromCsv(hasHeader bool) ([]product.Product, error) {

	file, err := os.Open(pl.csvPath)

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			log.Println("Failed to close file", file.Name())
		}
	}(file)

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	var headers []string
	if hasHeader {
		headers, err = reader.Read()
		if err != nil {
			return nil, err
		}
	}

	var products []product.Product
	pl.tagFieldMap = make(map[string]string)

	t := reflect.TypeOf(product.Product{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		pl.tagFieldMap[tag] = field.Name
	}

	for {
		line, err := reader.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("EOF")
				break
			}
			log.Println("Error", err)
			break
		}

		p := product.Product{}
		for i, value := range line {
			if value == "" {
				continue
			}
			pReflect := reflect.ValueOf(&p).Elem()
			if _, exists := pl.tagFieldMap[headers[i]]; !exists {
				continue
			}
			field := pReflect.FieldByName(pl.tagFieldMap[headers[i]])

			switch field.Type().String() {
			case "string":
				field.SetString(value)
			case "int":
				v, err := strconv.Atoi(value)
				if err != nil {
					log.Println("error converting")
				}
				field.SetInt(int64(v))
			case "float32":
				v, err := strconv.ParseFloat(value, 32)
				if err != nil {
					log.Println("error converting")
				}
				field.SetFloat(v)
			}
		}
		products = append(products, p)
	}
	return products, nil
}
