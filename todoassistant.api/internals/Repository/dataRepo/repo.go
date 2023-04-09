package dataRepo

import (
	"context"
	"test-va/internals/entity/dataEntity"
)

type DataRepository interface {
	GetCountries(ctx context.Context) ([]*dataEntity.Country, error)
}
