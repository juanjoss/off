package ports

import (
	"github.com/juanjoss/off-etl/model"
)

type Repository interface {
	/*
		Adds a product to the database.
	*/
	AddProduct(*model.Product) error

	/*
		Adds a brand model to the database.
	*/
	AddBrand(*model.Brand) error

	/*
		Gets a brand model from the database.
	*/
	GetBrand(string) (*model.Brand, error)

	/*
		Returns true if brands are already loaded, false otherwise.
	*/
	BrandsLoaded() bool

	/*
		Adds the brand models associated with a product barcode to the database.
	*/
	AddProductBrands(string, []*model.Brand) error

	/*
		Adds the specified nutrient levels model to the database and returns the inserted record's id.
	*/
	AddProductNutrientLevels(*model.NutrientLevels) (uint8, error)

	/*
		Gets the id from the specified nutrient levels model from the database.
	*/
	GetProductNutrientLevelsId(*model.NutrientLevels) (uint8, error)
}
