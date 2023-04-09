package dataService

import (
	"context"
	"log"
	datarepo "test-va/internals/Repository/dataRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/dataEntity"
	"time"
)

type DataService interface {
	GetCountries() ([]*dataEntity.Country, *ResponseEntity.ServiceError)
}

type dataSrv struct {
	repo datarepo.DataRepository
}

func (srv *dataSrv) GetCountries() ([]*dataEntity.Country, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	countries, err := srv.repo.GetCountries(ctx)
	if countries == nil {
		// log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return countries, nil

}

func NewDataService(repo datarepo.DataRepository) DataService {
	return &dataSrv{repo: repo}
}
