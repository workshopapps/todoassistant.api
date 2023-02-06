package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"test-va/internals/Repository/projectRepo"
	"test-va/internals/entity/projectEntity"
)

type sqlRepo struct{
	conn *sql.DB
}

func (s *sqlRepo) PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq)error{
	stmt := fmt.Sprintf(`INSERT INTO Projects(project_id, title, color, user_id, date_created)
						VALUES ('%v','%v','%v','%v','%v')`,
					req.ProjectId, req.Title, req.Color,req.UserId, req.CreatedAt)
	_, err := s.conn.Exec(stmt)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func NewProjectSqlRepo(conn *sql.DB) projectRepo.ProjectRepository{
	return &sqlRepo{conn: conn}
}
