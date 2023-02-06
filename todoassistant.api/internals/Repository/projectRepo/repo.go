package projectRepo

import (
	"context"
	"test-va/internals/entity/projectEntity"
)

type ProjectRepository interface {
	PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq) error
	GetListOfProjects(ctx context.Context, userId string) ([]*projectEntity.GetAllUserProjectRes, error)
}
