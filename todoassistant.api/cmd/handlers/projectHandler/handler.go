package projectHandler

import (
	"log"
	"net/http"

	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/projectEntity"
	"test-va/internals/service/projectService"

	"github.com/gin-gonic/gin"
)

type projectHandler struct {
	srv projectService.ProjectService
}

func NewProjectHandler(srv projectService.ProjectService) *projectHandler {
	return &projectHandler{srv: srv}
}

func (p *projectHandler) CreateProject(c *gin.Context) {
	var req projectEntity.CreateProjectReq
	value := c.GetString("userId")
	log.Println("userId is: ", value)
	if value == "" {

		c.AbortWithStatusJSON(http.StatusUnauthorized, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "you are not allowed to access this resource", nil, nil))
		return
	}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}
	req.UserId = value
	log.Println("create project req", req)

	project, errRes := p.srv.PersistProject(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error creating Project", errRes, nil))
		return
	}

	c.JSON(http.StatusOK,
		ResponseEntity.BuildSuccessResponse(http.StatusOK, "Created Project Successfully", project, nil))

}

func (p *projectHandler) GetAllUsersProjects(c *gin.Context) {
	userId := c.GetString("userId")
	log.Println(userId)
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "No userId found", nil, nil))
		return
	}
	projects, errRes := p.srv.GetListOfUsersProjects(userId)
	if projects == nil {
		message := "user with id " + userId + " has no project"
		c.AbortWithStatusJSON(http.StatusOK,
			ResponseEntity.BuildSuccessResponse(http.StatusNoContent, message, projects, nil))
		return
	}
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failure To Find all users project", errRes, nil))
		return
	}

	c.JSON(http.StatusOK,
		ResponseEntity.BuildSuccessResponse(http.StatusOK, "Users projects returned successfully", projects, nil))
}

func (p *projectHandler) EditProjectById(c *gin.Context) {
	var req projectEntity.EditProjectReq
	err := c.ShouldBind(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}

	projectId := c.Params.ByName("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no projectId id was provided", nil, nil))
		return
	}
	req.ProjectId = projectId

	userId := c.GetString("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Authentication Error, Invalid UserId", nil, nil))
		return
	}
	req.UserId = userId

	result, errRes := p.srv.EditProjectByID(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Unable to edit project", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Project edit successfully", result, nil)
	c.JSON(http.StatusOK, rd)
}

// Handle Delete task by id
func (p *projectHandler) DeleteProjectById(c *gin.Context) {
	projectId := c.Params.ByName("projectId")
	if projectId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no projectId id was provided", nil, nil))
		return
	}
	userId := c.GetString("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Authentication Error, Invalid UserId", nil, nil))
		return
	}
	_, errRes := p.srv.DeleteProjectByID(projectId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Unable to delete project", errRes, nil))
		return
	}
	rd := ResponseEntity.BuildSuccessResponse(200, "Project deleted successfully", nil, nil)
	c.JSON(http.StatusOK, rd)
}
