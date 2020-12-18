package dataloader

import (
	"time"

	"github.com/kiwisheets/gql-server/dataloader/generated"
	"github.com/kiwisheets/gql-server/orm/model"
	"gorm.io/gorm"
)

func newAddressByAddresseeIDLoader(db *gorm.DB, addresseeType string) *generated.AddressLoader {
	return generated.NewAddressLoader(generated.AddressLoaderConfig{
		MaxBatch: 1000,
		Wait:     1 * time.Millisecond,
		Fetch: func(addresseeIDs []int64) ([]*model.Address, []error) {
			addressRows, err := db.Model(&model.Address{}).Where("addressee_id IN (?) AND addressee_type = ?", addresseeIDs, addresseeType).Rows()

			if err != nil {
				if addressRows == nil {
					return nil, []error{err}
				}
			}
			defer addressRows.Close()

			addressByAddressee := map[int64]*model.Address{}

			for addressRows.Next() {
				var address model.Address
				db.ScanRows(addressRows, &address)
				if address.ID == 0 {
					// no value returned
				} else {
					addressByAddressee[int64(address.AddresseeID)] = &address
				}
			}

			orderedAddresses := make([]*model.Address, len(addresseeIDs))
			for i, id := range addresseeIDs {
				orderedAddresses[i] = addressByAddressee[id]
			}

			return orderedAddresses, nil
		},
	})
}
