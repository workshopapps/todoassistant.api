package mySqlRepo

import (
	"database/sql"
	"fmt"
	"log"

	"test-va/internals/Repository/reminderRepo"
	"test-va/internals/entity/reminderEntity"
	"test-va/internals/entity/taskEntity"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) CreateNewTask(req *taskEntity.CreateTaskReq) error {
	stmt := fmt.Sprintf(`INSERT INTO Tasks(
                  task_id,
                  user_id,
                  title,
                  description,
                  start_time,
                  end_time,
                  created_at,
                  va_option,
                  repeat_frequency
                  )
	VALUES ('%v','%v','%v','%v','%v','%v','%v','%v','%v')
	`, req.TaskId, req.UserId, req.Title, req.Description, req.StartTime, req.EndTime, req.CreatedAt, req.VAOption, req.Repeat)
	_, err := s.conn.Exec(stmt)
	if err != nil {
		log.Println(stmt)
		log.Println(err)
		return err
	}
	return nil
}

func (s *sqlRepo) SetTaskToExpired(id string) error {
	stmt := fmt.Sprintf(`UPDATE Tasks SET status = 'EXPIRED' WHERE task_id ='%v'`, id)
	_, err := s.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlRepo) GetAllUsersPendingTasks() ([]reminderEntity.GetPendingTasks, error) {
	stmt := `
		SELECT T.task_id, T.user_id, T.title,T.description, T.end_time, N.device_id
		FROM Tasks T join Notification_Tokens N on T.user_id = N.user_id
		WHERE status = 'PENDING';
	`

	var tasks []reminderEntity.GetPendingTasks
	query, err := s.conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		var task reminderEntity.GetPendingTasks
		var deviceId string
		err = query.Scan(&task.TaskId, &task.UserId, &task.Title, &task.Description, &task.EndTime, &deviceId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func NewSqlRepo(conn *sql.DB) reminderRepo.ReminderRepository {
	return &sqlRepo{conn: conn}
}
