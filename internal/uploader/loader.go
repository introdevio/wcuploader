package uploader

import (
	"errors"
	"github.com/introdevio/wcuploader/internal/product"
	"path/filepath"
	"regexp"
)

type ProductLoader struct {
	csvPath        string
	rootDir        string
	tagFieldMap    map[string]string
	pathRegex      *regexp.Regexp
	variationRegex *regexp.Regexp
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
	return pl.LoadFromPath()
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
		category := filepath.Base(filepath.Dir(file))
		if len(fields) != 4 {
			return nil, []error{errors.New("product does not follow right pattern")}
		}
		sku := fields[1]
		colorCode := fields[2]
		colorName := fields[3]
		img := product.NewLocalImageFromPath(file)

		v := product.Color{
			Sku:          colorCode,
			RegularPrice: "45",
			Name:         colorName,
			Image:        &img,
		}

		if p, exists := parentMap[sku]; exists {
			p.Images = append(p.Images, &img)
			p.Variations[colorCode] = v
			p.Colors[colorName] = true
		} else {
			a := product.Color{
				Sku:          colorCode,
				RegularPrice: "45",
				Name:         colorName,
				Image:        &img,
			}
			variations := make(map[string]product.Color)
			colors := make(map[string]bool)
			variations[colorCode] = a
			colors[colorName] = true
			p = &product.Product{
				Sku:          sku,
				ProductType:  "variable",
				Categories:   []string{category},
				RegularPrice: "45",
				Images:       []*product.LocalImage{&img},
				Colors:       colors,
				Variations:   variations,
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
