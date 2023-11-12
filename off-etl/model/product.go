package model

import (
	"strings"

	"golang.org/x/exp/slices"
)

var MandatoryStateTags = []string{
	"en:product-name-completed",
	"en:brands-completed",
	"en:quantity-completed",
	"en:front-photo-selected",
	"en:nutrition-facts-completed",
	"en:checked",
}

type ProductsRes struct {
	Count    int          `json:"count"`
	Products []ProductRes `json:"products"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type ProductRes struct {
	Barcode        string         `json:"code"`
	Name           string         `json:"product_name"`
	Quantity       string         `json:"quantity"`
	Brands         []string       `json:"brands_tags"`
	ImageUrl       string         `json:"image_url"`
	NutrientLevels NutrientLevels `json:"nutrient_levels"`
	Nutriments     Nutriments     `json:"nutriments"`
	NutriscoreData NutriscoreData `json:"nutriscore_data"`
	StateTags      []string       `json:"states_tags"`
}

type Product struct {
	Barcode          string `db:"barcode"`
	Name             string `db:"name"`
	Quantity         string `db:"quantity"`
	ImageUrl         string `db:"image_url"`
	Brands           []*Brand
	Energy100g       float32 `db:"energy_100g"`
	EnergyServing    float32 `db:"energy_serving"`
	NutrientLevelsId uint8   `db:"nutrient_levels_id"`
	NovaGroup        uint8   `db:"nova_group"`
	NutriscoreScore  int8    `db:"nutriscore_score"`
	NutriscoreGrade  string  `db:"nutriscore_grade"`
}

func (pr *ProductRes) ToModel() (*Product, error) {
	product := &Product{
		Barcode:         pr.Barcode,
		Name:            pr.Name,
		Quantity:        pr.Quantity,
		ImageUrl:        pr.ImageUrl,
		Brands:          []*Brand{},
		Energy100g:      pr.Nutriments.Energy100g,
		EnergyServing:   pr.Nutriments.EnergyServing,
		NovaGroup:       pr.Nutriments.NOVA,
		NutriscoreScore: pr.NutriscoreData.Score,
		NutriscoreGrade: strings.ToUpper(pr.NutriscoreData.Grade),
	}

	return product, nil
}

func (pr *ProductRes) HasMandatoryStateTags() bool {
	for _, st := range MandatoryStateTags {
		if !slices.Contains(pr.StateTags, st) {
			return false
		}
	}

	return true
}

func (p *Product) HasWrongFields() bool {
	return p.Barcode == "" ||
		p.Name == "" ||
		p.Quantity == "" ||
		p.ImageUrl == "" ||
		p.NovaGroup > 4 || p.NovaGroup < 1 ||
		p.NutriscoreScore < -15 || p.NutriscoreScore > 40 ||
		p.NutriscoreGrade == ""
}
