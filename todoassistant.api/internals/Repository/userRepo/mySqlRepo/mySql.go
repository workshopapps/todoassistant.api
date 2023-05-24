package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/userEntity"
	"time"
)

type mySql struct {
	conn *sql.DB
}

func NewMySqlUserRepo(conn *sql.DB) userRepo.UserRepository {
	return &mySql{conn: conn}
}

func (m *mySql) AssignVAToUser(user_id, va_id string) error {
	query := fmt.Sprintf(`
		SELECT user_id, virtual_assistant_id
		FROM Users
		WHERE user_id = '%s'
	`, user_id)
	var userId string
	var vaId string
	err := m.conn.QueryRowContext(context.Background(), query).Scan(&userId, &vaId)
	if err != nil {
		switch {
		case err.Error() == "sql: Scan error on column index 1, name \"virtual_assistant_id\": converting NULL to string is unsupported":
			vaId = ""
		default:
			return err
		}
	}
	if vaId != "" {
		return fmt.Errorf("user already has a VA")
	}

	query = fmt.Sprintf(`
		UPDATE Users SET virtual_assistant_id = '%s' WHERE user_id = '%s'
	`, va_id, user_id)

	_, err = m.conn.ExecContext(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *mySql) GetUsers(page int) ([]*userEntity.UsersRes, error) {
	var allUsers []*userEntity.UsersRes
	limit := 20
	offset := limit * (page - 1)
	query := fmt.Sprintf(`SELECT user_id, email, first_name, last_name, phone, date_of_birth, date_created
							FROM Users
							ORDER BY user_id
							LIMIT %d
							OFFSET %d`, limit, offset)

	ctx := context.Background()
	rows, err := m.conn.QueryContext(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	fmt.Println(rows.Next())
	// for rows.NextResultSet() {
	for rows.Next() {
		var user userEntity.UsersRes
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.DateOfBirth,
			&user.DateCreated,
		)

		if err != nil {
			return allUsers, err
		}

		allUsers = append(allUsers, &user)
	}
	// }

	if err = rows.Err(); err != nil {
		return allUsers, err
	}
	return allUsers, nil
}

func (m *mySql) GetByEmail(email string) (*userEntity.GetByEmailRes, error) {
	query := fmt.Sprintf(`
		SELECT user_id, email, password, first_name, last_name, phone, COALESCE(gender, ''), avatar,COALESCE(occupation, ''), COALESCE(country_id, 0)
		FROM Users
		WHERE email = '%s'
	`, email)
	var user userEntity.GetByEmailRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&user.UserId,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Gender,
		&user.Avatar,
		&user.Occupation,
		&user.CountryId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (m *mySql) GetById(user_id string) (*userEntity.GetByIdRes, error) {
	query := fmt.Sprintf(`
		SELECT user_id, password, email, first_name, last_name, phone, COALESCE(gender, ''), avatar
		FROM Users
		WHERE user_id = '%s'
	`, user_id)

	var user userEntity.GetByIdRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&user.UserId,
		&user.Password,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Gender,
		&user.Avatar,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (m *mySql) Persist(req *userEntity.CreateUserReq) error {
	log.Println("from persist", req)
	ctx := context.Background()
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := fmt.Sprintf(` INSERT INTO Users(
                   user_id,
                   first_name,
                   last_name,
                   email,
                   phone,
                   password,
                   account_status
                   ) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v')`,
		req.UserId, req.FirstName, req.LastName, req.Email, req.Phone, req.Password, req.AccountStatus)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	err = m.userSettings(ctx, tx, req.UserId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Create User Settings
func (m *mySql) userSettings(ctx context.Context, tx *sql.Tx, userId string) error {
	notiResult, err := m.notificationSettings(ctx, tx, userId)
	if err != nil {
		return err
	}
	log.Println("Passed notification")

	nId, err := notiResult.LastInsertId()
	if err != nil {
		return err
	}

	prodResult, err := m.productEmailSettings(ctx, tx, userId)
	if err != nil {
		return err
	}

	pId, err := prodResult.LastInsertId()
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(`INSERT INTO User_Settings(
		user_id,
		notification_settings_id,
    	product_email_settings_id
	) VALUES('%v', '%v', '%v')`, userId, nId, pId)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

// Create Notification Settings
func (m *mySql) notificationSettings(ctx context.Context, tx *sql.Tx, userId string) (sql.Result, error) {
	stmt := fmt.Sprintf(`INSERT INTO Notification_Settings(
		user_id
	) VALUES('%v')`, userId)

	return tx.ExecContext(ctx, stmt)
}

// Create Product Email Settings
func (m *mySql) productEmailSettings(ctx context.Context, tx *sql.Tx, userId string) (sql.Result, error) {
	stmt := fmt.Sprintf(`INSERT INTO Product_Email_Settings(
		user_id
	) VALUES('%v')`, userId)

	return tx.ExecContext(ctx, stmt)
}

// Create function to update user in database
func (m *mySql) UpdateUser(req *userEntity.UpdateUserReq, userId string) error {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	stmt := fmt.Sprintf(`UPDATE Users SET
                 first_name ='%s',
                 last_name='%s',
                 email ='%s',
                 phone='%s',
                 gender='%s',
                 date_of_birth='%s',
				 occupation='%s',
				 country_id='%d' WHERE user_id ='%s'
                 `, req.FirstName, req.LastName, req.Email, req.Phone, req.Gender, req.DateOfBirth, req.Occupation, req.CountryId, userId)

	_, err := m.conn.ExecContext(ctx, stmt)
	log.Println("from repo", err)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) UpdateImage(userId, fileName string) error {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	stmt := fmt.Sprintf(`UPDATE Users SET
                 avatar = '%s'
                 WHERE user_id ='%s'
                 `, fileName, userId)

	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

// Auxillary function to update user
// func updateField(tx *sql.Tx, userId string, field string, val interface{}) (sql.Result, error) {
// 	return tx.Exec(fmt.Sprintf(`UPDATE Users SET %s = '%v' WHERE user_id = '%v'`, field, val, userId))
// }

// // Auxillary function to update user
// func updateFieldIfSet(tx *sql.Tx, userId string, field string, val interface{}) (sql.Result, error) {
// 	v, ok := val.(string)
// 	if ok && v != "" {
// 		return updateField(tx, userId, field, v)
// 	}
// 	return nil, nil
// }

func (m *mySql) ChangePassword(user_id, newPassword string) error {
	query := fmt.Sprintf(`UPDATE Users SET password = '%v' WHERE user_id = '%v'`, newPassword, user_id)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *mySql) DeleteUser(user_id string) error {
	query := fmt.Sprintf(`DELETE FROM Users WHERE user_id = "%s"`, user_id)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *mySql) AddToken(req *userEntity.ResetPasswordRes) error {
	stmt := fmt.Sprintf(` INSERT INTO Reset_Token(
                   token_id,
                   user_id,
                   token,
                   expiry
                   ) VALUES ('%v', '%v', '%v', '%v')`,
		req.TokenId, req.UserId, req.Token, req.Expiry)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) GetTokenById(token, userId string) (*userEntity.ResetPasswordWithTokenRes, error) {
	query := fmt.Sprintf(`
		SELECT token_id, user_id, token, expiry
		FROM Reset_Token
		WHERE token = '%s'
		AND user_id = '%s'
	`, token, userId)

	var tokenRes userEntity.ResetPasswordWithTokenRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&tokenRes.TokenId,
		&tokenRes.UserId,
		&tokenRes.Token,
		&tokenRes.Expiry,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tokenRes, nil
}

func (m *mySql) DeleteToken(userId string) error {
	query := fmt.Sprintf(`DELETE FROM Reset_Token WHERE user_id = "%s"`, userId)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// get user notification settings
func (m *mySql) GetNotificationSettingsById(userId string) (*userEntity.NotificationSettingsRes, error) {
	query := fmt.Sprintf(`
		SELECT IF(new_comments, 'true', 'false') as new_comments,
			IF(expired_tasks, 'true', 'false') as expired_task,
			IF(reminder_tasks, 'true', 'false') as reminder_task,
			IF(va_accepting_task, 'true', 'false') as va_accepting_task,
			IF(tasks_assigned_va, 'true', 'false') as task_assigned_va,
			IF(subscription, 'true', 'false') as subscription
		FROM Notification_Settings
		WHERE user_id = '%s'
	`, userId)
	var notificationSettings userEntity.NotificationSettingsRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&notificationSettings.NewComments,
		&notificationSettings.ExpiredTasks,
		&notificationSettings.ReminderTasks,
		&notificationSettings.VaAcceptingTask,
		&notificationSettings.TaskAssingnedVa,
		&notificationSettings.Subscribtion,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &notificationSettings, nil
}

// get product email settings
func (m *mySql) GetProductEmailSettingsById(userId string) (*userEntity.ProductEmailSettingsRes, error) {
	query := fmt.Sprintf(`
		SELECT IF(new_products, 'true', 'false') as new_products,
			IF(login_alert, 'true', 'false') as login_alert,
			IF(promotions_and_offers, 'true', 'false') as promotions_and_offers,
			IF(tips_daily_digest, 'true', 'false') as tips_daily_digest
		FROM Product_Email_Settings
		WHERE user_id = '%s'
	`, userId)
	var productEmailSettings userEntity.ProductEmailSettingsRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&productEmailSettings.NewProducts,
		&productEmailSettings.LoginAlert,
		&productEmailSettings.PromotionAndOffers,
		&productEmailSettings.TipsDailyDigest,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &productEmailSettings, nil
}

//set reminder settings

func (m *mySql) SetReminderSettings(req *userEntity.ReminderSettngsReq) error {
	ctx := context.Background()
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if a record already exists for the given userID
	checkStmt := fmt.Sprintf("SELECT COUNT(*) FROM Reminder_Settings WHERE user_id = '%v'", req.UserId)
	var count int
	err = tx.QueryRowContext(ctx, checkStmt).Scan(&count)
	if err != nil {
		return err
	}

	var stmt string
	if count > 0 {
		// Update the existing record
		stmt = fmt.Sprintf(`UPDATE Reminder_Settings SET
			remindMeVia = '%v',
			whenSnooze = '%v',
			autoReminder = '%v',
			reminderTime = '%v',
			refresh = '%v'
			WHERE user_id = '%v'`,
			req.RemindMeVia, req.WhenSnooze, req.AutoReminder, req.ReminderTime, req.Refresh, req.UserId)
	} else {
		// Insert a new record
		stmt = fmt.Sprintf(`INSERT INTO Reminder_Settings(
			remindMeVia,
			whenSnooze,
			autoReminder,
			reminderTime,
			refresh,
			user_id
		) VALUES ('%v', '%v', '%v', '%v', '%v', '%v')`,
			req.RemindMeVia, req.WhenSnooze, req.AutoReminder, req.ReminderTime, req.Refresh, req.UserId)
	}

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// get reminder settings for a user
func (m *mySql) GetReminderSettings(userId string) (*userEntity.ReminderSettngsRes, error) {
	stmt := fmt.Sprintf(`
		SELECT remindMeVia, whenSnooze, autoReminder, reminderTime, refresh
		FROM Reminder_Settings
		WHERE user_id = '%s'
	`, userId)
	var reminderSettings userEntity.ReminderSettngsRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, stmt).Scan(
		&reminderSettings.RemindMeVia,
		&reminderSettings.WhenSnooze,
		&reminderSettings.AutoReminder,
		&reminderSettings.ReminderTime,
		&reminderSettings.Refresh,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &reminderSettings, nil
}
