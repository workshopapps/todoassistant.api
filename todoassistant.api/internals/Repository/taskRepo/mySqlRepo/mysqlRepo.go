package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
	"test-va/internals/service/timeSrv"
	"time"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) AssignTaskToVa(ctx context.Context, vaId, taskId string) error {
	log.Println(vaId)
	log.Println(taskId)
	stmt := fmt.Sprintf(`UPDATE Tasks SET va_id ='%v' WHERE task_id ='%v'`, vaId, taskId)
	_, err := s.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlRepo) GetVADetails(ctx context.Context, userId string) (string, error) {
	var vaId *string
	stmt := fmt.Sprintf(`
SELECT
	virtual_Assistant_id from Users
WHERE user_id = '%v'
`, userId)
	row := s.conn.QueryRowContext(ctx, stmt)
	err := row.Scan(&vaId)
	if err != nil {
		return "", err
	}
	return *vaId, nil
}

func (s *sqlRepo) GetAllTaskAssignedToVA(ctx context.Context, vaId string) ([]*vaEntity.VATask, error) {
	stmt := fmt.Sprintf(`SELECT
    T.task_id,
    T.title,
    T.end_time,
    T.status,
    T.description,
	T.va_option,
	T.va_id,
    concat(U2.first_name, ' ', U2.last_name) AS 'name',
    T.user_id,
    U2.phone
FROM Tasks T
         join va_table U on T.va_id = U.va_id join Users U2 on U2.user_id = T.user_id
WHERE T.va_id = '%s'
ORDER BY T.created_at DESC
;`, vaId)

	queryRow, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var Results []*vaEntity.VATask

	for queryRow.Next() {
		var res vaEntity.VATask
		err := queryRow.Scan(&res.TaskId, &res.Title, &res.EndTime, &res.Status, &res.Description, &res.VaOption, &res.VaId, &res.User.Name, &res.User.UserId, &res.User.Phone)
		if err != nil {
			return nil, err
		}
		Results = append(Results, &res)
	}

	return Results, nil
}

// get all task and user details for VA
func (s *sqlRepo) GetAllTaskForVA(ctx context.Context) ([]*vaEntity.VATaskAll, error) {
	stmt := `SELECT
    T.task_id,
    T.title,
    T.end_time,
    T.status,
    T.description,
	T.va_option,
	T.comment_count,
	COALESCE(T.va_id, ''),
    concat(U.first_name, ' ', U.last_name) AS 'name',
    T.user_id,
    U.phone,
	U.avatar
	FROM Tasks T
        join  Users U on T.user_id = U.user_id
	ORDER BY T.created_at DESC;`

	queryRow, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var Results []*vaEntity.VATaskAll

	for queryRow.Next() {
		var res vaEntity.VATaskAll
		err := queryRow.Scan(&res.TaskId, &res.Title, &res.EndTime, &res.Status, &res.Description, &res.VaOption, &res.CommentCount, &res.VaId, &res.User.Name, &res.User.UserId, &res.User.Phone, &res.User.Avatar)
		if err != nil {
			return nil, err
		}
		Results = append(Results, &res)
	}

	return Results, nil
}

func (s *sqlRepo) GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error) {
	query := fmt.Sprintf(`
		SELECT task_id, user_id, title, description, start_time, end_time, status
		FROM Tasks
		WHERE user_id = '%s' AND status = 'PENDING'
	`, userId)

	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*taskEntity.GetPendingTasksRes

	for rows.Next() {
		var task taskEntity.GetPendingTasksRes
		err := rows.Scan(
			&task.TaskId,
			&task.UserId,
			&task.Title,
			&task.Description,
			&task.StartTime,
			&task.EndTime,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *sqlRepo) PersistAndAssign(ctx context.Context, req *taskEntity.CreateTaskReq) error {
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

	var vaId string
	stmt3 := fmt.Sprintf(`
		SELECT
			virtual_Assistant_id from Users
		WHERE user_id = '%v'
		`, req.UserId)
	row := tx.QueryRowContext(ctx, stmt3)
	err = row.Scan(&vaId)
	if err != nil {
		log.Println("3", err)
		return err
	}

	stmt := fmt.Sprintf(`INSERT
		INTO Tasks(
					task_id,
                  user_id,
                  title,
                  description,
                  start_time,
                  end_time,
                  created_at,
                  va_option,
                  repeat_frequency,
		           va_id
				   )
		VALUES ('%v','%v','%v','%v','%v','%v','%v', '%v', '%v', '%v')`, req.TaskId, req.UserId, req.Title, req.Description,
		req.StartTime, req.EndTime, req.CreatedAt, req.VAOption, req.Repeat, vaId)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		log.Println("1", err)
		return err
	}

	for _, file := range req.Files {
		stmt2 := fmt.Sprintf(`INSERT
		INTO Taskfiles(
		               task_id,
		               file_link,
		               file_type
		               )
		VALUES ('%v', '%v', '%v')`, req.TaskId, file.FileLink, file.FileType)
		_, err = tx.ExecContext(ctx, stmt2)
		if err != nil {
			log.Println("2", err)
			return err
		}
	}

	return nil
}

func (s *sqlRepo) Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()
	log.Println("create task req", req)

	stmt := fmt.Sprintf(`INSERT
		INTO Tasks(
				task_id,
				user_id,
				title,
				description,
				start_time,
				end_time,
				created_at,
				va_option,
				repeat_frequency,
				notify,
				project_id,
				scheduled_date
			)
		VALUES ('%v','%v','%v','%v','%v','%v','%v', '%v', '%v',%t, '%v', '%v')`, req.TaskId, req.UserId, req.Title, req.Description,
		req.StartTime, req.EndTime, req.CreatedAt, req.VAOption, req.Repeat, req.Notify, req.ProjectId, req.ScheduledDate)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, file := range req.Files {
		stmt2 := fmt.Sprintf(`INSERT
								INTO Taskfiles(
								task_id,	
								file_link,
								file_type
							)
							VALUES ('%v', '%v', '%v')`, req.TaskId, file.FileLink, file.FileType)
		_, err = tx.ExecContext(ctx, stmt2)
		if err != nil {
			log.Println("err", err)
			return err
		}
	}

	return nil
}

// search by name
func (s *sqlRepo) SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	// defer func() {
	// 	if err != nil {
	// 		tx.Rollback()
	// 	} else {
	// 		tx.Commit()
	// 	}
	// }()

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, start_time
		FROM Tasks
		WHERE title LIKE '%s%%'
	`, title.SearchQuery)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Searchedtasks []*taskEntity.SearchTaskRes

	for rows.Next() {
		var singleTask taskEntity.SearchTaskRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			&singleTask.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		Searchedtasks = append(Searchedtasks, &singleTask)
	}
	return Searchedtasks, nil
}

// get task by ID

func (s *sqlRepo) GetTaskByID(ctx context.Context, taskId string) (*taskEntity.GetTasksByIdRes, error) {
	var task taskEntity.GetTasksByIdRes
	tim := timeSrv.NewTimeStruct()
	tx, err := s.conn.BeginTx(ctx, nil)
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

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, description, status, start_time, repeat_frequency, end_time, created_at, COALESCE(updated_at, ""), COALESCE(va_id,""), notify, COALESCE(project_id,""), COALESCE(scheduled_date,"")
		FROM Tasks T
		WHERE task_id = '%s'`, taskId)

	row := tx.QueryRow(stmt)
	if err := row.Scan(
		&task.TaskId,
		&task.UserId,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.StartTime,
		&task.Repeat,
		&task.EndTime,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.VaId,
		&task.Notify,
		&task.ProjectId,
		&task.ScheduledDate,
	); err != nil {
		return nil, err
	}

	var features taskEntity.TaskFeatures
	if task.VaId != "" {
		features.IsAssigned = true
	}
	if task.ScheduledDate != "" {
		features.IsScheduled = true
	}

	end, err := time.Parse(time.RFC3339, task.EndTime)
	if err != nil {
		return nil, err
	}

	if tim.TimeBefore(end) && task.Status == "PENDING" {
		log.Println(tim.TimeBefore(end))
		features.IsExpired = true
	}

	if task.Status == "COMPLETED" {
		features.IsCompleted = true
	}

	task.TaskFeatures = features

	stmt2 := fmt.Sprintf(`
		SELECT F.file_link, F.file_type
		FROM Tasks AS T
		JOIN Taskfiles as F
		ON T.task_id = F.task_id
		WHERE F.task_id = '%s'
	`, taskId)

	log.Println("Created AT", task)
	rows, err := tx.QueryContext(ctx, stmt2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var taskFile taskEntity.TaskFile

		err := rows.Scan(
			&taskFile.FileLink,
			&taskFile.FileType,
		)
		if err != nil {
			return nil, err
		}
		task.Files = append(task.Files, taskFile)
	}

	return &task, nil
}

func (s *sqlRepo) GetListOfExpiredTasks(ctx context.Context) ([]*taskEntity.GetAllExpiredRes, error) {
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT task_id, user_id, title, start_time
				FROM Tasks
				WHERE status = 'EXPIRED'`

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Searchedtasks []*taskEntity.GetAllExpiredRes

	for rows.Next() {
		var singleTask taskEntity.GetAllExpiredRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			&singleTask.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		Searchedtasks = append(Searchedtasks, &singleTask)
	}
	return Searchedtasks, nil
}

func (s *sqlRepo) GetListOfPendingTasks(ctx context.Context) ([]*taskEntity.GetAllPendingRes, error) {
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT task_id, user_id, title, end_time
				FROM Tasks
				WHERE status = 'PENDING'`

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var SearchedPendingtasks []*taskEntity.GetAllPendingRes

	for rows.Next() {
		var singleTask taskEntity.GetAllPendingRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			// &singleTask.VAOption,
			&singleTask.EndTime,
		)
		// fmt.Println(err)
		if err != nil {
			return nil, err
		}
		SearchedPendingtasks = append(SearchedPendingtasks, &singleTask)
	}
	return SearchedPendingtasks, nil
}

// Get All task
func (s *sqlRepo) GetAllTasks(ctx context.Context, userId string) ([]*taskEntity.GetAllTaskRes, error) {
	tim := timeSrv.NewTimeStruct()
	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}
	log.Println("HERE ", userId)
	stmt := fmt.Sprintf(`
		SELECT task_id, title, description, status, start_time, repeat_frequency, end_time, created_at, COALESCE(updated_at, ""), COALESCE(va_id,""), notify, COALESCE(project_id,""), COALESCE(scheduled_date,"")
		FROM Tasks T WHERE user_id = '%s'`, userId)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AllTasks []*taskEntity.GetAllTaskRes

	for rows.Next() {
		var singleTask taskEntity.GetAllTaskRes

		if err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.Title,
			&singleTask.Description,
			&singleTask.Status,
			&singleTask.StartTime,
			&singleTask.Repeat,
			&singleTask.EndTime,
			&singleTask.CreatedAt,
			&singleTask.UpdatedAt,
			&singleTask.VaId,
			&singleTask.Notify,
			&singleTask.ProjectId,
			&singleTask.ScheduledDate,
		); err != nil {
			log.Println("error ", err)
			return nil, err
		}

		var features taskEntity.TaskFeatures
		if singleTask.VaId != "" {
			features.IsAssigned = true
		}
		if singleTask.ScheduledDate != "" {
			features.IsScheduled = true
		}

		end, err := time.Parse(time.RFC3339, singleTask.EndTime)
		if err != nil {
			return nil, err
		}

		if tim.TimeBefore(end) && singleTask.Status == "PENDING" {
			features.IsExpired = true
		}

		if singleTask.Status == "COMPLETED" {
			features.IsCompleted = true
		}

		singleTask.TaskFeatures = features
		AllTasks = append(AllTasks, &singleTask)
	}
	return AllTasks, nil
}

// Delete task by id
func (s *sqlRepo) DeleteTaskByID(ctx context.Context, taskId string) error {

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
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`Delete from Comments Where task_id = '%s'`, taskId))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`Delete from Tasks  WHERE task_id = '%s'`, taskId))
	if err != nil {
		return err
	}
	return nil
}

// Delete All
func (s *sqlRepo) DeleteAllTask(ctx context.Context, userId string) error {

	//var res taskEntity.GetTasksByIdRes
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
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`Delete from Tasks where  user_id = '%s'`, userId))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *sqlRepo) EditTaskById(ctx context.Context, taskId string, req *taskEntity.EditTaskReq) error {
	notifyInt := 0
	if req.Notify {
		notifyInt = 1
	}

	stmt := fmt.Sprintf(`UPDATE Tasks SET
							title = '%s',
							description = '%s',
							status = '%s',
							start_time = '%s',
							repeat_frequency = '%s',
							end_time = '%s',
							updated_at = '%s',
							notify = '%d',
							project_id ='%s',
							scheduled_date= '%s'
							WHERE task_id = '%s'
						`, req.Title, req.Description, req.Status, req.StartTime, req.Repeat, req.EndTime, req.UpdatedAt, notifyInt, req.ProjectId, req.ScheduledDate, taskId)

	log.Println(req.ProjectId)
	_, err := s.conn.ExecContext(ctx, stmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *sqlRepo) UpdateTaskStatusByID(ctx context.Context, taskId string, req *taskEntity.UpdateTaskStatus) error {
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
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`UPDATE Tasks SET status = '%s' WHERE task_id = '%s'`, req.Status, taskId))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// comment
func (s *sqlRepo) PersistComment(ctx context.Context, req *taskEntity.CreateCommentReq) error {
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

	stmt := fmt.Sprintf(`INSERT INTO Comments(
                  sender_id,
                  task_id,
                  comment,
				  created_at,
				  status,
				  isEmoji
                  )
	VALUES ('%v','%v','%v','%v','%v','%v')
	`, req.SenderId, req.TaskId, req.Comment, req.CreatedAt, req.Status, req.IsEmoji)

	stmt2 := fmt.Sprintf(`UPDATE Tasks SET comment_count=comment_count+1 WHERE task_id='%v'`, req.TaskId)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, stmt2)
	if err != nil {
		return err
	}
	// _, err := s.conn.Exec(stmt)
	if err != nil {
		log.Println(stmt)
		log.Println(err)
		return err
	}

	return nil
}

// get all comments
func (s *sqlRepo) GetAllComments(ctx context.Context, taskId string) ([]*taskEntity.GetCommentRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(`
		SELECT id, sender_id, task_id, comment, created_at,status,isEmoji
		FROM Comments WHERE task_id = '%s'`, taskId)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AllComment []*taskEntity.GetCommentRes

	for rows.Next() {
		var singleTask taskEntity.GetCommentRes

		err := rows.Scan(
			&singleTask.Id,
			&singleTask.SenderId,
			&singleTask.TaskId,
			&singleTask.Comment,
			&singleTask.CreatedAt,
			&singleTask.Status,
			&singleTask.IsEmoji,
		)
		if err != nil {
			return nil, err
		}
		AllComment = append(AllComment, &singleTask)
	}
	return AllComment, nil
}

// get all comments
func (s *sqlRepo) GetComments(ctx context.Context) ([]*taskEntity.GetCommentRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT id, sender_id, task_id, comment, created_at, status, isEmoji FROM Comments`

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AllComment []*taskEntity.GetCommentRes

	for rows.Next() {
		var singleTask taskEntity.GetCommentRes

		err := rows.Scan(
			&singleTask.Id,
			&singleTask.SenderId,
			&singleTask.TaskId,
			&singleTask.Comment,
			&singleTask.CreatedAt,
			&singleTask.Status,
			&singleTask.IsEmoji,
		)
		if err != nil {
			return nil, err
		}
		AllComment = append(AllComment, &singleTask)
	}
	return AllComment, nil
}

// Delete comment by id
func (s *sqlRepo) DeleteCommentByID(ctx context.Context, commentId string) error {
	log.Println("hererer", commentId)
	_, err := s.conn.ExecContext(ctx, fmt.Sprintf(`Delete from Comments  WHERE id = '%s'`, commentId))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func NewSqlRepo(conn *sql.DB) taskRepo.TaskRepository {
	return &sqlRepo{conn: conn}
}
