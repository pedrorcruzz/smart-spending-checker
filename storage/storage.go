package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pedrorcruzz/smart-spending-checker/product"
)

var dataFile = filepath.Join("data", "products.json")

func LoadProducts() (product.ProductList, error) {
	var list product.ProductList
	file, err := os.Open(dataFile)
	if err != nil {
		return product.ProductList{
			Products:      []product.Product{},
			MonthlyProfit: 0,
			Month:         int(time.Now().Month()),
			Year:          time.Now().Year(),
		}, nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&list)
	return list, err
}

func SaveProducts(list product.ProductList) error {
	os.MkdirAll(filepath.Dir(dataFile), 0755)
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(list)
}
