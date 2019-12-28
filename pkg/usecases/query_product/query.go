package query_product

import (
	"github.com/webronin26/snap-up/pkg/entities"
	"github.com/webronin26/snap-up/pkg/store"
)

func Exec(productID int32) (entities.Product, error) {

	var product entities.Product

	query := store.DB.Model(entities.Product{}).
		Where("id = ?", productID).
		Scan(&product)

	if err := query.Error; err != nil {
		return product, err
	}

	return product, nil
}
