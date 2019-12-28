package add_record

import (
	"time"

	"github.com/webronin26/snap-up/pkg/entities"
	"github.com/webronin26/snap-up/pkg/store"
)

type Input struct {
	ProductID     int32
	ProductNumber int32
	CustomerID    int32
	Number        int32
}

func Exec(input Input) error {

	record := new(entities.Record)
	record.Number = input.Number
	record.CustomerID = input.CustomerID
	record.ProductID = input.ProductID
	record.Date = time.Now()

	add := store.DB.Model(entities.Record{}).Create(&record)
	if err := add.Error; err != nil {
		return err
	}

	return nil
}
