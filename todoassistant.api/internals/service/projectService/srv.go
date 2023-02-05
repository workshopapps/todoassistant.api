package projectService

import (
	"context"
	"log"
	"test-va/internals/Repository/projectRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/projectEntity"
	"test-va/internals/service/loggerService"
	"test-va/internals/service/timeSrv"
	"test-va/internals/service/validationService"
	"time"
)

type ProjectService interface {
	PersistProject(req *projectEntity.CreateProjectReq)(*projectEntity.CreateProjectRes, *ResponseEntity.ServiceError)
}

type projectSrv struct {
	repo 			projectRepo.ProjectRepository
	timeSrv       	timeSrv.TimeService
	validationSrv 	validationService.ValidationSrv
	logger			loggerService.LogSrv
}

func NewProjectSrv(repo projectRepo.ProjectRepository, timeSrv timeSrv.TimeService, validationSrv validationService.ValidationSrv, logger loggerService.LogSrv) ProjectService{
	return &projectSrv{repo: repo, timeSrv: timeSrv, validationSrv: validationSrv, logger: logger}
}

func (p *projectSrv) PersistProject(req *projectEntity.CreateProjectReq)(*projectEntity.CreateProjectRes, *ResponseEntity.ServiceError){
	_, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	err := p.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}
	return nil,nil
}
