package update_product

import (
	"github.com/webronin26/snap-up/pkg/entities"
	"github.com/webronin26/snap-up/pkg/store"
)

type Input struct {
	ProductID  int32
	ItemNumber int32
}

func Exec(input Input) error {

	update := store.DB.Model(entities.Product{}).
		Where("id = ?", input.ProductID).
		Update("item_number", input.ItemNumber)

	if err := update.Error; err != nil {
		return err
	}

	return nil
}
