package routes

import (
	"test-va/cmd/handlers/dataHandler"

	"test-va/internals/service/dataService"

	"github.com/gin-gonic/gin"
)

func DataRoutes(v1 *gin.RouterGroup, service dataService.DataService) {

	handler := dataHandler.NewDataHandler(service)
	project := v1.Group("/data")

	project.GET("/countries", handler.GetCountries)

}
