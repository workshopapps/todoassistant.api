package userRepo

import (
	"test-va/internals/entity/userEntity"
)

type UserRepository interface {
	GetUsers(page int) ([]*userEntity.UsersRes, error)
	Persist(req *userEntity.CreateUserReq) error
	GetByEmail(email string) (*userEntity.GetByEmailRes, error)
	GetById(user_id string) (*userEntity.GetByIdRes, error)
	UpdateUser(req *userEntity.UpdateUserReq, userId string) error
	DeleteUser(user_id string) error
	ChangePassword(userId, newPassword string) error
}
