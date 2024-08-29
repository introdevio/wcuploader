package uploader

import (
	"errors"
	"github.com/introdevio/wcuploader/internal"
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

func (pl *ProductLoader) Load() ([]internal.Product, []error) {
	return pl.LoadFromPath()
}

func (pl *ProductLoader) LoadFromPath() ([]internal.Product, []error) {
	files, err := filepath.Glob(filepath.Join(pl.rootDir, "*", "*.jpg"))

	if err != nil {
		return nil, []error{err}
	}
	parentMap := make(map[string]*internal.Product)

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
		img := internal.NewLocalImageFromPath(file)

		v := internal.Color{
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
			a := internal.Color{
				Sku:          colorCode,
				RegularPrice: "45",
				Name:         colorName,
				Image:        &img,
			}
			variations := make(map[string]internal.Color)
			colors := make(map[string]bool)
			variations[colorCode] = a
			colors[colorName] = true
			p = &internal.Product{
				Sku:          sku,
				ProductType:  "variable",
				Categories:   []string{category},
				RegularPrice: "45",
				Images:       []*internal.LocalImage{&img},
				Colors:       colors,
				Variations:   variations,
			}
			parentMap[sku] = p
		}
	}

	// load media

	var products []internal.Product
	for _, value := range parentMap {
		products = append(products, *value)
	}
	return products, nil
}
