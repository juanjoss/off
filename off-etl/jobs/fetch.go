package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/juanjoss/off-etl/model"
)

const (
	countryDomain = "world"
)

var (
	/*
		Product fields used as parameters for a products request URL.
	*/
	productFields = []string{
		"code",
		"product_name",
		"quantity",
		"brands_tags",
		"image_url",
		"nutrient_levels",
		"nutriments",
		"nutriscore_grade",
		"states_tags",
	}
)

/*
	Fetches products from the Open Foods Facts API.

	pageSize indicates the number of products to be fetched per page.

	page indicates the page number to fetch.
*/
func FetchProducts(page, pageSize int) (*model.ProductsRes, error) {
	url := fmt.Sprintf(
		"https://%s.openfoodfacts.org/api/v2/search?%s&json=true&pageSize=%d&page=%d",
		countryDomain,
		strings.Join(productFields, ","),
		pageSize,
		page,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching products (status is not 200)")
	}

	var products model.ProductsRes
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}

/*
	Fetches brands from the Open Foods Facts API.
*/
func FetchBrands() (*model.BrandsRes, error) {
	url := fmt.Sprintf("https://%s.openfoodfacts.org/brands.json", countryDomain)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching brands (status is not 200)")
	}

	var brands model.BrandsRes
	err = json.NewDecoder(res.Body).Decode(&brands)
	if err != nil {
		return nil, err
	}

	return &brands, nil
}
