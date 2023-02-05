package mySqlRepo

import (
	"context"
	"database/sql"
	"test-va/internals/Repository/projectRepo"
	"test-va/internals/entity/projectEntity"
)

type sqlRepo struct{
	conn *sql.DB
}

func (s *sqlRepo) PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq)error{

	return nil
}

func NewProjectSqlRepo(conn *sql.DB) projectRepo.ProjectRepository{
	return &sqlRepo{conn: conn}
}
