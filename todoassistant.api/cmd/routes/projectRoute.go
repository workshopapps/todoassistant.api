package routes

import (
	"test-va/cmd/handlers/projectHandler"
	"test-va/cmd/middlewares"

	"test-va/internals/service/projectService"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(v1 *gin.RouterGroup, service projectService.ProjectService, srv tokenservice.TokenSrv) {

	jwtMWare := middlewares.NewJWTMiddleWare(srv)

	handler := projectHandler.NewProjectHandler(service)
	project := v1.Group("/project")


	project.Use(jwtMWare.ValidateJWT())
	{
		project.POST("", handler.CreateProject)
		project.GET("/", handler.GetAllUsersProjects)
	}


}
