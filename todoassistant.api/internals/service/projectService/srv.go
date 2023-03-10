package projectService

import (
	"context"
	"log"
	"net/http"
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
	PersistProject(req *projectEntity.CreateProjectReq) (*projectEntity.CreateProjectRes, *ResponseEntity.ServiceError)
	GetListOfUsersProjects(userId string) ([]*projectEntity.GetProjectRes, *ResponseEntity.ServiceError)
	EditProjectByID(req *projectEntity.EditProjectReq) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
	DeleteProjectByID(projectId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError)
}

type projectSrv struct {
	repo          projectRepo.ProjectRepository
	timeSrv       timeSrv.TimeService
	validationSrv validationService.ValidationSrv
	logger        loggerService.LogSrv
}

func (p *projectSrv) PersistProject(req *projectEntity.CreateProjectReq) (*projectEntity.CreateProjectRes, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	err := p.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}
	//set created time
	req.CreatedAt = p.timeSrv.CurrentTimeString() //.Format(time.RFC3339)
	req.ProjectId = uuid.New().String()

	err = p.repo.PersistProject(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	data := projectEntity.CreateProjectRes{
		ProjectId: req.ProjectId,
		UserId:    req.UserId,
		Title:     req.Title,
		Color:     req.Color,
	}
	return &data, nil
}

func (p *projectSrv) GetListOfUsersProjects(userId string) ([]*projectEntity.GetProjectRes, *ResponseEntity.ServiceError) {

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
	return projects, nil
}

func (p *projectSrv) EditProjectByID(req *projectEntity.EditProjectReq) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	err := p.validationSrv.Validate(req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewValidatingError("Bad Data Input")
	}

	project, err := p.repo.GetProject(ctx, req.ProjectId, req.UserId)
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	req = Check(req, project.Color, project.Title)
	req.UpdatedAt = p.timeSrv.CurrentTime().Format(time.RFC3339)

	result, err := p.repo.EditProject(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "Project updated successfully", result, nil), nil
}

func (p *projectSrv) DeleteProjectByID(projectId string) (*ResponseEntity.ResponseMessage, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	err := p.repo.DeleteProjectByID(ctx, projectId)
	if err != nil {
		log.Println(err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	return ResponseEntity.BuildSuccessResponse(http.StatusOK, "project Deleted successfully", nil, nil), nil
}

func NewProjectSrv(repo projectRepo.ProjectRepository, timeSrv timeSrv.TimeService, validationSrv validationService.ValidationSrv, logger loggerService.LogSrv) ProjectService {
	return &projectSrv{repo: repo, timeSrv: timeSrv, validationSrv: validationSrv, logger: logger}
}

func Check(req *projectEntity.EditProjectReq, color, title string) *projectEntity.EditProjectReq {
	if req.Color == "" {
		req.Color = color
	}
	if req.Title == "" {
		req.Title = title
	}

	return req
}
