package vaHandler

import (
	"context"
	"log"
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/tokenEntity"
	"test-va/internals/entity/vaEntity"
	"test-va/internals/service/taskService"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/userService"
	"test-va/internals/service/vaService"

	"github.com/gin-gonic/gin"
)

type vaHandler struct {
	tokenSrv tokenservice.TokenSrv
	vaSrv    vaService.VAService
	taskSrv  taskService.TaskService
	userSrv  userService.UserSrv
}

func NewVaHandler(tokenSrv tokenservice.TokenSrv, vaSrv vaService.VAService, taskSrv taskService.TaskService, userSrv userService.UserSrv) *vaHandler {
	return &vaHandler{tokenSrv: tokenSrv, vaSrv: vaSrv, taskSrv: taskSrv, userSrv: userSrv}
}

func (v *vaHandler) UpdateVA(c *gin.Context) {
	param := c.Param("va_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No id in url"})
	}
	var req vaEntity.EditVaReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest,
				"Bad Input Data", err, nil))
		return
	}

	user, errRes := v.vaSrv.UpdateVA(&req, param)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Authorization Error", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Changed user details Successfully", user, nil))
}

func (v *vaHandler) DeleteVA(c *gin.Context) {
	param := c.Param("va_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No id in url"})
	}

	errRes := v.vaSrv.DeleteVA(param)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Error Deleting Message", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"delete user successful", nil, nil))

}

func (v *vaHandler) ChangePassword(c *gin.Context) {
	var req vaEntity.ChangeVAPassword

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest,
				"Bad Input Data", err, nil))
		return
	}

	errRes := v.vaSrv.ChangePassword(&req)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Authorization Error", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Changed Password Successful", nil, nil))
}

func (v *vaHandler) FindByEmail(c *gin.Context) {
	param := c.Param("email")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No email in url"})
	}

	user, errRes := v.vaSrv.FindByEmail(param)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Authorization Error", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Found User By Email Successful", user, nil))
}

func (v *vaHandler) SignUp(c *gin.Context) {
	var req vaEntity.CreateVAReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest,
				"Bad Input Data", err, nil))
		return
	}

	user, serviceError := v.vaSrv.SignUp(&req)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Failed to Sign Up", serviceError, nil))
		return
	}

	token, s, err := v.tokenSrv.CreateToken(user.VaId, user.AccountType, user.Email)
	if err != nil {
		return
	}

	// handle token service
	var tokenData = &tokenEntity.TokenRes{
		Token:        token,
		RefreshToken: s,
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusCreated,
		"created user successfully", user, tokenData))
}

func (v *vaHandler) Login(c *gin.Context) {
	var req vaEntity.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest,
				"Bad Input Data", err, nil))
		return
	}

	user, serviceError := v.vaSrv.Login(&req)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Authorization Error", serviceError, nil))
		return
	}

	token, s, err := v.tokenSrv.CreateToken(user.VaId, user.AccountType, user.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failed to create token", nil, nil))
		return
	}

	// handle token service
	var tokenData = &tokenEntity.TokenRes{
		Token:        token,
		RefreshToken: s,
	}
	ctx := c.Request.Context()
	ctxII := c.Request.WithContext(context.WithValue(ctx, "id", "testId"))
	c.Request = ctxII
	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusCreated,
		"Login user successful", user, tokenData))
}

func (v *vaHandler) GetVAByID(c *gin.Context) {
	vaId := c.Param("va_id")
	if vaId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No Id in url"})
	}

	va, errRes := v.vaSrv.GetVA(vaId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Authorization Error", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "Found User By Id Successful", va, nil))
}

func (v *vaHandler) GetUserAssignedToVA(c *gin.Context) {
	param := c.Param("va_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No Id in url"})
		return
	}

	va, serviceError := v.vaSrv.GetAllUserToVa(param)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Error getting Users", serviceError, nil))
		return
	}

	if len(va) == 0 {
		c.JSON(http.StatusOK, ResponseEntity.BuildErrorResponse(http.StatusOK,
			"No User Found", "", nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Found Users Successfully", va, nil))

}

func (v *vaHandler) GetTaskByUser(c *gin.Context) {
	param := c.Param("user_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No Id in url"})
		return
	}
	log.Println(param)

	task, serviceError := v.taskSrv.GetAllTask(param)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Error Task ", serviceError, nil))
		return
	}

	if len(task) == 0 {
		c.JSON(http.StatusOK, ResponseEntity.BuildErrorResponse(http.StatusOK,
			"No Task Found", "", nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Found Tasks Successfully", task, nil))
}

func (v *vaHandler) GetSingleUserProfile(c *gin.Context) {
	param := c.Param("user_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No Id in url"})
		return
	}
	log.Println(param)

	user, serviceError := v.userSrv.GetUser(param)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Error getting user ", serviceError, nil))
		return
	}

	// if len(task) == 0 {
	// 	c.JSON(http.StatusOK, ResponseEntity.BuildErrorResponse(http.StatusOK,
	// 		"No Task Found", "", nil))
	// 	return
	// }

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Found User Successfully", user, nil))
}

func (v *vaHandler) GetAllAssignedUsersTask(c *gin.Context) {
	param := c.Param("va_id")
	if param == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": "No Va Id in url"})
		return
	}
	log.Println(param)

	task, serviceError := v.taskSrv.GetTaskAssignedToVA(param)
	if serviceError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusInternalServerError,
				"Error Task ", serviceError, nil))
		return
	}

	if len(task) == 0 {
		c.JSON(http.StatusOK, ResponseEntity.BuildErrorResponse(http.StatusOK,
			"No Task Found", "", nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK,
		"Found Tasks Successfully", task, nil))
}
