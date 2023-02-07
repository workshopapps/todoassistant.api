package mysqlRepo

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

func (s *sqlRepo) GetListOfProjects(ctx context.Context, userId string) ([]*projectEntity.GetAllUserProjectRes, error){
	stmt := fmt.Sprintf(`
		SELECT project_id, title, color, user_id FROM Projects
		WHERE user_id = '%s'
	`, userId)

	rows, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*projectEntity.GetAllUserProjectRes

	for rows.Next() {
		var project projectEntity.GetAllUserProjectRes
		err := rows.Scan(
			&project.ProjectId,
			&project.Title,
			&project.Color,
			&project.UserId,

		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	if rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func NewProjectSqlRepo(conn *sql.DB) projectRepo.ProjectRepository{
	return &sqlRepo{conn: conn}
}
