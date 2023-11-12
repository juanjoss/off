package jobs

import (
	"log"

	"github.com/juanjoss/off-etl/model"
	"github.com/juanjoss/off-etl/ports"
)

var (
	apiPageNumber = 1
	apiPageSize   = 100
)

func RunProductsETL(repo ports.Repository) {
	load(repo, transform(repo, extract(apiPageNumber)))
	apiPageNumber++

	log.Printf("products (page=%d) load process finished", apiPageNumber-1)
}

func extract(page int) <-chan model.ProductRes {
	responses := make(chan model.ProductRes)

	productsRes, err := FetchProducts(page, apiPageSize)
	if err != nil {
		log.Fatalf("error fetching products: %v", err)
	}

	go func() {
		defer close(responses)
		for _, p := range productsRes.Products {
			responses <- p
		}
	}()

	return responses
}

func transform(repo ports.Repository, productResponses <-chan model.ProductRes) <-chan *model.Product {
	transformedProducts := make(chan *model.Product)

	go func() {
		defer close(transformedProducts)

		for pr := range productResponses {
			/*
				check mandatory state tags to avoid dirty data.
			*/
			if !pr.HasMandatoryStateTags() {
				continue
			}

			/*
				convert product response to product model
			*/
			productModel, err := pr.ToModel()
			if err != nil {
				continue
			}

			/*
				check for wrong fields in product model
			*/
			if productModel.HasWrongFields() {
				continue
			}

			/*
				search if nutrient levels exist, if not, add them to the database.
			*/
			id, err := repo.GetProductNutrientLevelsId(&pr.NutrientLevels)
			if err != nil {
				id, err = repo.AddProductNutrientLevels(&pr.NutrientLevels)
				if err != nil {
					log.Printf("error loading nutrient levels for product %v: %v", pr.Name, err)
				}
			}

			if id == 0 {
				continue
			}

			productModel.NutrientLevelsId = id

			/*
				search product brands and load them.
			*/
			for _, brandTag := range pr.Brands {
				brand, err := repo.GetBrand(brandTag)
				if err != nil {
					continue
				}

				productModel.Brands = append(productModel.Brands, brand)
			}

			/*
				add the product model for the load phase.
			*/
			transformedProducts <- productModel
		}
	}()

	return transformedProducts
}

func load(repo ports.Repository, productModels <-chan *model.Product) {
	for {
		productModel, ok := <-productModels
		if ok {
			/*
				load the product model.
			*/
			err := repo.AddProduct(productModel)
			if err != nil {
				log.Printf("error loading product model %v: %v", productModel, err)
			}

			/*
				load product brands.
			*/
			err = repo.AddProductBrands(productModel.Barcode, productModel.Brands)
			if err != nil {
				log.Printf("error loading product brands: %v", err.Error())
			}
		} else {
			return
		}
	}
}
