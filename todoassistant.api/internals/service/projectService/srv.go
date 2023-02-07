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

	"github.com/google/uuid"
)

type ProjectService interface {
	PersistProject(req *projectEntity.CreateProjectReq)(*projectEntity.CreateProjectRes, *ResponseEntity.ServiceError)
	GetListOfUsersProjects(userId string) ([]*projectEntity.GetAllUserProjectRes, *ResponseEntity.ServiceError)
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
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	err := p.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}
	//set created time
	req.CreatedAt = p.timeSrv.CurrentTime().Format(time.RFC3339)
	req.ProjectId = uuid.New().String()

	err = p.repo.PersistProject(ctx, req)
	if err != nil {
			log.Println(err)
			return nil, ResponseEntity.NewInternalServiceError(err)
	}

	data := projectEntity.CreateProjectRes{
		ProjectId: req.ProjectId,
		UserId: req.UserId,
		Title: req.Title,
		Color: req.Color,
	}
	return &data,nil
}

func (p *projectSrv) GetListOfUsersProjects(userId string) ([]*projectEntity.GetAllUserProjectRes, *ResponseEntity.ServiceError){

	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	projects, err := p.repo.GetListOfProjects(ctx, userId)
	if projects == nil {
		// log.Println("no rows returned")
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return projects,nil
}
