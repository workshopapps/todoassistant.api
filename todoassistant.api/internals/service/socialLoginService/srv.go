package socialLoginService

import (
	"strings"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/timeSrv"
	tokenservice "test-va/internals/service/tokenService"
	"time"

	"github.com/google/uuid"
)

type LoginSrv interface {
	LoginResponse(req *userEntity.GoogleLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError)
	FacebookLoginResponse(req *userEntity.FacebookLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError)
}

type loginSrv struct {
	repo     userRepo.UserRepository
	timeSrv  timeSrv.TimeService
	tokenSrv tokenservice.TokenSrv
}

// Google login godoc
// @Summary	Login user using google account
// @Description	Google login route
// @Tags	Social Login
// @Accept	json
// @Produce	json
// @Param	request	body	userEntity.GoogleLoginReq	true "Google login"
// @Success	200  {object}  userEntity.LoginRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/googlelogin [post]
func (l *loginSrv) LoginResponse(req *userEntity.GoogleLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError) {
	user, _ := l.repo.GetByEmail(req.Email)
	if user == nil {

		resData := &userEntity.CreateUserReq{
			UserId:        uuid.New().String(),
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Email:         req.Email,
			AccountStatus: "ACTIVE",
			DateCreated:   l.timeSrv.CurrentTime().Format(time.RFC3339),
		}

		err := l.repo.Persist(resData)

		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	}

	user, _ = l.repo.GetByEmail(req.Email)
	// tokenSrv := tokenservice.NewTokenSrv("fvmvmvmvf")

	accessToken, refreshToken, err := l.tokenSrv.CreateToken(user.Email, "user", user.UserId)

	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	loginUser := &userEntity.LoginRes{
		UserId:       user.UserId,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Gender:       user.Gender,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	return loginUser, nil
}

// Facebook login godoc
// @Summary	Login user using facebook account
// @Description	Facebook login route
// @Tags	Social Login
// @Accept	json
// @Produce	json
// @Param	request	body	userEntity.FacebookLoginReq	true "Facebook login"
// @Success	200  {object}  userEntity.LoginRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/facebooklogin [post]
func (l *loginSrv) FacebookLoginResponse(req *userEntity.FacebookLoginReq) (*userEntity.LoginRes, *ResponseEntity.ServiceError) {
	user, _ := l.repo.GetByEmail(req.Email)
	name := strings.Split(req.Name, " ")

	firstName := name[0]
	lastName := name[1]

	if user == nil {
		resData := &userEntity.CreateUserReq{
			UserId:        uuid.New().String(),
			FirstName:     firstName,
			LastName:      lastName,
			Email:         req.Email,
			AccountStatus: "ACTIVE",
			DateCreated:   l.timeSrv.CurrentTime().Format(time.RFC3339),
		}

		err := l.repo.Persist(resData)

		if err != nil {
			return nil, ResponseEntity.NewInternalServiceError(err)
		}
	}

	user, _ = l.repo.GetByEmail(req.Email)
	// tokenSrv := tokenservice.NewTokenSrv("fvmvmvmvf")

	accessToken, refreshToken, err := l.tokenSrv.CreateToken(user.Email, "user", user.UserId)

	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	loginUser := &userEntity.LoginRes{
		UserId:       user.UserId,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Gender:       user.Gender,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	return loginUser, nil
}

func NewLoginSrv(repo userRepo.UserRepository, timeSrv timeSrv.TimeService, tokenSrv tokenservice.TokenSrv) LoginSrv {
	return &loginSrv{repo, timeSrv, tokenSrv}
}
