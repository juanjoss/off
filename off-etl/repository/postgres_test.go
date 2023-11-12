package repository

import (
	"strings"
	"testing"

	"github.com/juanjoss/off-etl/model"
)

var pr *postgresRepo

func Init() {
	pr = NewRepository()
}

func TestPostgresRepo(t *testing.T) {
	Init()

	// add a brand
	brand := &model.Brand{
		Tag: "coca-cola",
	}
	err := pr.AddBrand(brand)
	if err != nil {
		t.Fatalf("error inserting brand: %v", err.Error())
	}

	// create product nutritional information
	nutrientLevels := model.NutrientLevels{
		Fat:          "low",
		SaturatedFat: "low",
		Sugar:        "low",
		Salt:         "low",
	}
	nutrientLevelsId, err := pr.AddProductNutrientLevels(&nutrientLevels)
	if err != nil {
		t.Fatalf("error while getting nutrient levels id: %v", err.Error())
	}

	nutriscoreData := model.NutriscoreData{
		Grade: "b",
		Score: 1,
	}

	nutriments := model.Nutriments{
		Energy100g:    1,
		EnergyServing: 3.3,
		NOVA:          4,
	}

	// add a product
	p := &model.Product{
		Barcode:         "5449000131805",
		Name:            "Coca-Cola Zero",
		Quantity:        "0.33 l",
		ImageUrl:        "https://images.openfoodfacts.org/images/products/544/900/013/1805/front_en.490.400.jpg",
		Energy100g:      nutriments.Energy100g,
		EnergyServing:   nutriments.EnergyServing,
		NovaGroup:       nutriments.NOVA,
		NutriscoreScore: nutriscoreData.Score,
		NutriscoreGrade: strings.ToUpper(nutriscoreData.Grade),
	}
	p.NutrientLevelsId = nutrientLevelsId

	err = pr.AddProduct(p)
	if err != nil {
		t.Fatalf("error inserting product: %v", err.Error())
	}

	// add product brands
	var brands []*model.Brand
	brands = append(brands, brand)
	err = pr.AddProductBrands(p.Barcode, brands)
	if err != nil {
		t.Fatalf("error inserting product brands: %v", err.Error())
	}
}
