package uploader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/introdevio/wcuploader/internal/product"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type ProductLoader struct {
	csvPath     string
	rootDir     string
	tagFieldMap map[string]string
}

func NewCSVProductLoader(csvPath, rootDir string) *ProductLoader {
	return &ProductLoader{
		csvPath: csvPath,
		rootDir: rootDir,
	}
}

func NewPathProductLoader(rootDir string) *ProductLoader {
	return &ProductLoader{
		rootDir: rootDir,
	}
}

func (pl *ProductLoader) Load() ([]product.Product, []error) {
	if pl.csvPath != "" {
		return pl.loadFromCsv(true)
	} else {
		return pl.LoadFromPath()
	}
}

func (pl *ProductLoader) loadFromCsv(hasHeader bool) ([]product.Product, []error) {

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

		st := mapLineToProduct(line, headers, errs)
		st.LoadMedia(pl.rootDir)
		items = append(items, st)
	}
	return items, errs
}

func (pl *ProductLoader) LoadFromPath() ([]product.Product, []error) {
	files, err := filepath.Glob(filepath.Join(pl.rootDir, "*", "*.jpg"))

	if err != nil {
		return nil, []error{err}
	}
	parentMap := make(map[string]product.Product)

	// create parents or products with no variations
	for _, file := range files {
		category := filepath.Base(filepath.Dir(file))
		fileName := filepath.Base(file)
		productFields := strings.Split(fileName, "-")
		var sku string
		var p product.Product

		fullSku := strings.ReplaceAll(productFields[0], " ", "-")
		fullVariation := strings.Split(fullSku, "-")
		if len(fullVariation) > 0 {
			sku = fullVariation[0]
			p = product.Product{
				PostTitle:    fullSku,
				Sku:          fullSku,
				Stock:        1,
				StockStatus:  "instock",
				Backorders:   "no",
				RegularPrice: 45,
			}
			if parent, exists := parentMap[sku]; exists {
				parent.Children = append(parent.Children, p)
				continue
			} else {
				parent = product.Product{
					PostTitle:   sku,
					Sku:         sku,
					StockStatus: "instock",
					ProductType: "variable",
					Stock:       1,
					Categories:  []string{category},
				}
				parent.Children = []product.Product{p}
				parentMap[sku] = parent
			}

		} else {
			sku = fullSku
			p = product.Product{
				PostTitle:    sku,
				Tags:         nil,
				Categories:   []string{category},
				Sku:          sku,
				Stock:        1,
				StockStatus:  "instock",
				Backorders:   "no",
				RegularPrice: 45,
			}
			parentMap[sku] = p
		}
	}

	// load media

	for _, file := range files {
		fileName := filepath.Base(file)
		productFields := strings.Split(fileName, "-")

		sku := strings.ReplaceAll(productFields[0], " ", "-")
		fullVariation := strings.Split(sku, "-")

		if len(fullVariation) > 0 {
			p := parentMap[fullVariation[0]]
			children := p.Children
			for i, _ := range children {
				c := &children[i]
				if c.Sku == sku {
					c.Images = append(c.Images, file)
				}
			}
		} else {
			p := parentMap[fullVariation[0]]
			p.Images = append(p.Images, file)
		}
	}

	for _, value := range parentMap {
		fmt.Printf("%+v\n\n\n", value)
	}
	return nil, nil
}

func mapLineToProduct(line []string, headers []string, errs []error) product.Product {
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

	return st
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
