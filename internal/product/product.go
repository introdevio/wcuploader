package product

type Product struct {
	Description      string
	ShortDescription string
	ProductType      string
	Categories       []string
	Sku              string
	RegularPrice     string
	SalePrice        string
	Images           []*LocalImage
	Variations       map[string]Color
	Colors           map[string]bool
}

type Color struct {
	Id           string
	Name         string
	Description  string
	Sku          string
	RegularPrice string
	SalePrice    string
	Image        *LocalImage
}

type LocalImage struct {
	Path          string
	RemoteImageId int
	RemoteUrl     string
}

func NewLocalImageFromPath(path string) LocalImage {
	return LocalImage{
		Path: path,
	}
}
