package dataHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/service/dataService"

	"github.com/gin-gonic/gin"
)

type dataHandler struct {
	srv dataService.DataService
}

func NewDataHandler(srv dataService.DataService) *dataHandler {
	return &dataHandler{srv: srv}
}

func (d *dataHandler) GetCountries(c *gin.Context) {
	response, err := d.srv.GetCountries()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Error getting countries data", err, nil))
		return
	}
	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "Countries data", response, nil))
}
