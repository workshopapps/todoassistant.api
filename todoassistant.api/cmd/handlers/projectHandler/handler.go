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

func (p *projectHandler) CreateProject(c *gin.Context){
	var req projectEntity.CreateProjectReq
	value := c.GetString("userId")
	log.Println("userId is: ", value)
	if value == "" {
		log.Println("112")
		c.AbortWithStatusJSON(http.StatusUnauthorized, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "you are not allowed to access this resource", nil, nil))
		return
	}
	log.Println("create project req",req)
}


