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
	"regexp"
	"strconv"
	"strings"
)

var CategoryMap = map[string]int{"damas": 22, "caballeros": 23}

type ProductLoader struct {
	csvPath        string
	rootDir        string
	tagFieldMap    map[string]string
	pathRegex      *regexp.Regexp
	variationRegex *regexp.Regexp
}

func NewCSVProductLoader(csvPath, rootDir string) *ProductLoader {
	return &ProductLoader{
		csvPath: csvPath,
		rootDir: rootDir,
	}
}

func NewPathProductLoader(rootDir string) *ProductLoader {
	filePathRegex, err := regexp.Compile(`/(\w+)(?:-|\s)+(\w+)(?:-|\s)+(\w+)(?:-|\s).+$`)
	if err != nil {
		panic(err)
	}

	return &ProductLoader{
		pathRegex: filePathRegex,
		rootDir:   rootDir,
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
	parentMap := make(map[string]*product.Product)

	// create parents or products with no variations
	for _, file := range files {
		fields := pl.pathRegex.FindStringSubmatch(file)
		category := filepath.Base(file)
		if len(fields) != 4 {
			return nil, []error{errors.New("product does not follow right pattern")}
		}
		sku := fields[1]
		colorCode := fields[2]
		colorName := fields[3]

		v := product.Variation{
			Sku:          sku + "-" + colorCode,
			RegularPrice: "45",
			Status:       "publish",
			ManageStock:  true,
			StockStatus:  "instock",
			Attribute:    product.VariationAttribute{Option: colorName},
		}

		if p, exists := parentMap[sku]; exists {
			p.Children = append(p.Children, v)
			color := *p.Attributes[0]
			color.Options = append(color.Options, colorName)
		} else {
			a := product.Attribute{
				Name:      "Color",
				Visible:   true,
				Variation: true,
				Options:   []string{colorName},
			}
			p = &product.Product{
				PostTitle:    sku,
				ProductType:  "variable",
				Categories:   []product.Category{{Id: CategoryMap[category]}},
				Sku:          sku,
				Stock:        true,
				StockStatus:  "instock",
				Status:       "publish",
				RegularPrice: "45",
				Images:       nil,
				Children:     []product.Variation{v},
				Attributes:   []*product.Attribute{&a},
			}
			parentMap[sku] = p
		}
	}

	// load media

	var products []product.Product
	for _, value := range parentMap {
		products = append(products, *value)
	}
	return products, nil
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
