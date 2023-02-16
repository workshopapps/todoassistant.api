package projectRepo

import (
	"context"
	"test-va/internals/entity/projectEntity"
)

type ProjectRepository interface {
	PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq) error
	GetListOfProjects(ctx context.Context, userId string) ([]*projectEntity.GetProjectRes, error)
	GetProject(ctx context.Context, projectId, userId string) (*projectEntity.GetProjectRes, error)
	EditProject(ctx context.Context, req *projectEntity.EditProjectReq) (*projectEntity.EditProjectRes, error)
	DeleteProjectByID(ctx context.Context, projectId string) error
}
