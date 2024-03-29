package reminderService

import (
	"fmt"
	"log"
	"os"
	notiRepo "test-va/internals/Repository/notificationRepo/mysqlRepo"
	"test-va/internals/Repository/reminderRepo/mySqlRepo"
	"test-va/internals/data-store/mysql"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/service/notificationService"
	"test-va/internals/service/validationService"
	"testing"
	"time"

	firebase "firebase.google.com/go"
	"github.com/go-co-op/gocron"
)

func Test_reminderSrv_SetReminder(t *testing.T) {
	taskId := "fd077d56-763b-43ac-9f0e-5fbe6f30cbc1"
	dsn := os.Getenv("dsn")
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()
	repo := mySqlRepo.NewSqlRepo(conn)
	notificationRepo := notiRepo.NewMySqlNotificationRepo(conn)
	validationSrv := validationService.NewValidationStruct()

	notificationSrv := notificationService.New(&firebase.App{}, notificationRepo, validationSrv)
	if err != nil {
		fmt.Println("Could Not Send Message", err)
	}

	gcrn := gocron.NewScheduler(time.UTC)
	srv := NewReminderSrv(gcrn, repo, notificationSrv)

	var data taskEntity.CreateTaskReq
	due := time.Now().Add(15 * time.Second).Format(time.RFC3339)
	data.EndTime = due
	data.TaskId = taskId
	srv.SetReminder(&data)

	time.Sleep(5 * time.Minute)
}

func Test_reminderSrv_SetReminderEveryXMin(t *testing.T) {
	//dsn := os.Getenv("dsn")
	//if dsn == "" {
	//	dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	//}
	//
	//connection, err := mysql.NewMySQLServer(dsn)
	//if err != nil {
	//	log.Println("Error Connecting to DB: ", err)
	//	return
	//}
	//defer connection.Close()
	//conn := connection.GetConn()
	//conn.Ping()
	//fmt.Println(time.Now().Format(time.RFC3339), time.Minute*2)

	//gcrn := gocron.NewScheduler(time.UTC)
	//srv := NewReminderSrv(gcrn, conn)
	//
	//srv.SetReminderEveryXMin(30)

	t2 := time.Now().Add(5 * time.Minute * 1420)

	hours := time.Until(t2).Hours() / 24
	fmt.Println(hours)

}

func Test_Empty(t *testing.T) {
	s := gocron.NewScheduler(time.UTC)
	tt := time.Now().Add(time.Second * 10)
	log.Println("starting", time.Now().Format(time.Kitchen))
	s.Every(5).StartAt(tt).Do(func() {
		log.Println("Doing something")
		log.Println("finishing", time.Now().Format(time.Kitchen))
	})
	s.LimitRunsTo(1)
	s.StartAsync()
	time.Sleep(2 * time.Minute)
}
