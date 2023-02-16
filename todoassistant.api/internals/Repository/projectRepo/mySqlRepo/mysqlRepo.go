package mySqlRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"test-va/internals/Repository/projectRepo"
	"test-va/internals/entity/projectEntity"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) PersistProject(ctx context.Context, req *projectEntity.CreateProjectReq) error {
	stmt := fmt.Sprintf(`INSERT INTO Projects(project_id, title, color, user_id, date_created)
						VALUES ('%v','%v','%v','%v','%v')`,
		req.ProjectId, req.Title, req.Color, req.UserId, req.CreatedAt)
	_, err := s.conn.Exec(stmt)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *sqlRepo) GetListOfProjects(ctx context.Context, userId string) ([]*projectEntity.GetProjectRes, error) {
	stmt := fmt.Sprintf(`
		SELECT project_id, title, color, user_id FROM Projects
		WHERE user_id = '%s'
	`, userId)

	rows, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*projectEntity.GetProjectRes

	for rows.Next() {
		var project projectEntity.GetProjectRes
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

func (s *sqlRepo) GetProject(ctx context.Context, projectId, userId string) (*projectEntity.GetProjectRes, error) {
	stmt := fmt.Sprintf(`
		SELECT project_id, title, color, user_id FROM Projects
		WHERE project_id = '%s' AND user_id = '%s'
	`, projectId, userId)

	row := s.conn.QueryRowContext(ctx, stmt)
	if row == nil {
		return nil, errors.New("no project with that id")
	}

	var project projectEntity.GetProjectRes
	err := row.Scan(
		&project.ProjectId,
		&project.Title,
		&project.Color,
		&project.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Delete project by id
func (s *sqlRepo) DeleteProjectByID(ctx context.Context, projectId string) error {

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`DELETE from Tasks WHERE project_id = '%s'`, projectId))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`DELETE from Projects WHERE project_id = '%s'`, projectId))
	if err != nil {
		return err
	}
	return nil
}

func (m *sqlRepo) EditProject(ctx context.Context, req *projectEntity.EditProjectReq) (*projectEntity.EditProjectRes, error) {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := fmt.Sprintf(`UPDATE Projects 
						SET
						title = '%v',
						color = '%v',
						date_updated = '%v'
						WHERE user_id = '%v' AND project_id = '%v'`,
		req.Title, req.Color, req.UpdatedAt, req.UserId, req.ProjectId)

	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &projectEntity.EditProjectRes{
		ProjectId: req.ProjectId,
		UserId:    req.UserId,
		Title:     req.Title,
		Color:     req.Color,
		UpdatedAt: req.UpdatedAt,
	}, nil
}

func NewProjectSqlRepo(conn *sql.DB) projectRepo.ProjectRepository {
	return &sqlRepo{conn: conn}
}
