package projectRepo

import (
	"context"
	"test-va/internals/entity/projectEntity"
)

type ProjectRepository interface {
	PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq) error
}
