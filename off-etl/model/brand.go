package model

type BrandsRes struct {
	Brands []BrandRes `json:"tags"`
}

type BrandRes struct {
	Tag string `json:"id"`
}

type Brand struct {
	Tag string `db:"tag"`
}

func (br *BrandRes) ToModel() *Brand {
	return &Brand{
		Tag: br.Tag,
	}
}
