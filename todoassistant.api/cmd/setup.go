package cmd

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"test-va/cmd/handlers/paymentHandler"
	"test-va/cmd/middlewares"
	"test-va/cmd/routes"
	mySqlCallRepo "test-va/internals/Repository/callRepo/mySqlRepo"
	mySqlRepo5 "test-va/internals/Repository/dataRepo/mySqlRepo"
	mySqlNotifRepo "test-va/internals/Repository/notificationRepo/mysqlRepo"
	projectMysqlRepo "test-va/internals/Repository/projectRepo/mySqlRepo"
	mySqlRemindRepo "test-va/internals/Repository/reminderRepo/mySqlRepo"
	mySqlRepo4 "test-va/internals/Repository/subscribeRepo/mySqlRepo"
	"test-va/internals/Repository/taskRepo/mySqlRepo"
	mySqlRepo2 "test-va/internals/Repository/userRepo/mySqlRepo"
	mySqlRepo3 "test-va/internals/Repository/vaRepo/mySqlRepo"
	awss3 "test-va/internals/amazon/awsS3"
	"test-va/internals/data-store/mysql"
	firebaseinit "test-va/internals/firebase-init"
	"test-va/internals/msg-queue/Emitter"
	"test-va/internals/service/awsService"
	"test-va/internals/service/callService"
	"test-va/internals/service/cryptoService"
	"test-va/internals/service/dataService"
	"test-va/internals/service/emailService"
	log_4_go "test-va/internals/service/loggerService/log-4-go"
	"test-va/internals/service/notificationService"
	"test-va/internals/service/projectService"
	"test-va/internals/service/reminderService"
	"test-va/internals/service/socialLoginService"
	"test-va/internals/service/subscribeService"
	"test-va/internals/service/taskService"
	"test-va/internals/service/timeSrv"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/userService"
	"test-va/internals/service/vaService"
	"test-va/internals/service/validationService"
	"test-va/utils"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"

	"github.com/getsentry/sentry-go"

	"github.com/go-co-op/gocron"

	_ "test-va/docs"

	"github.com/stripe/stripe-go/v74"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() {
	// set up sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://cffbc6e0ff4c480bb9ad07108811485b@o4504281294176256.ingest.sentry.io/4504282768539648",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// try to connect to rabbit mq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	emitter, err := Emitter.NewAmqpEmitter(rabbitConn)
	if err != nil {
		log.Println("error:", err)
		return
	}

	//Load configurations
	config, err := utils.LoadConfig("./")

	//load google config file
	utils.LoadGoogleConfig()

	if err != nil {
		log.Fatal("cannot load config file", err)
	}

	stripe.Key = config.StripeKey

	if config.StripeKey == "" {

		stripe.Key = "sk_test_51M9xknFf5hgzULIC40q0q9nzGz6ByBYNrFYzgUB2zsVfDZwhhiss5fi3OmLVhzOwxLfnT4bMqjj9Uh4oaLQrCRhU00EUIT0yl3"

	}

	dsn := config.DataSourceName
	if dsn == "" {
		dsn = "hawaiian_comrade:YfqvJUSF43DtmH#^ad(K+pMI&@(team-ruler-todo.c6qozbcvfqxv.ap-south-1.rds.amazonaws.com:3306)/todoDB"
	}

	port := config.SeverAddress
	if port == "" {
		port = "2022"
	}

	secret := config.TokenSecret
	if secret == "" {
		log.Fatal("secret key not found")
	}

	fromEmailAddr := config.FromEmailAddr
	if fromEmailAddr == "" {
		log.Fatal("smtp email sender address not found")
	}

	smtpPWD := config.SMTPpwd
	if smtpPWD == "" {
		log.Fatal("smtp password not found")
	}

	smtpHost := config.SMTPhost
	if fromEmailAddr == "" {
		log.Fatal("smtp host address not found")
	}

	smtpPort := config.SMTPport
	if fromEmailAddr == "" {
		log.Fatal("smtp port not found")
	}

	AWSAccess := config.AWSAccess
	if AWSAccess == "" {
		log.Fatal("AWS Access Key not found")
	}

	AWSSecret := config.AWSSecret
	if AWSSecret == "" {
		log.Fatal("AWS Secret Key not found")
	}

	//Repo

	//db service
	connection, err := mysql.NewMySQLServer(dsn)
	if err != nil {
		log.Println("Error Connecting to DB: ", err)
		return
	}
	defer connection.Close()
	conn := connection.GetConn()

	projectRepo := projectMysqlRepo.NewProjectSqlRepo(conn)
	// task repo service
	taskRepo := mySqlRepo.NewSqlRepo(conn)

	//user repo service
	userRepo := mySqlRepo2.NewMySqlUserRepo(conn)

	//notification repo service
	notificationRepo := mySqlNotifRepo.NewMySqlNotificationRepo(conn)

	//reminder repo service
	remindRepo := mySqlRemindRepo.NewSqlRepo(conn)

	//va repo service
	vaRepo := mySqlRepo3.NewVASqlRepo(conn)

	// subscribe repo
	subRepo := mySqlRepo4.NewMySqlSubscribeRepo(conn)

	// data repo
	dataRepo := mySqlRepo5.NewDataSqlRepo(conn)

	//SERVICES

	//time service
	timeSrv := timeSrv.NewTimeStruct()

	//validation service
	validationSrv := validationService.NewValidationStruct()

	//Notification Service
	//Note Handle Unable to Connect to Firebase
	firebaseApp, err := firebaseinit.SetupFirebase()
	if err != nil {
		fmt.Println("UNABLE TO CONNECT TO FIREBASE", err)
	}
	notificationSrv := notificationService.New(firebaseApp, notificationRepo, validationSrv)
	if err != nil {
		fmt.Println("Could Not Send Message", err)
	}

	// s3 init
	s3session, err := awss3.NewAWSSession(AWSAccess, AWSSecret, "")
	if err != nil {
		log.Println("Error Connecting to AWS S3: ", err)
		return
	}
	// fmt.Println(s3session)
	// create cron tasks for checking if time is due

	callRepo := mySqlCallRepo.NewSqlCallRepo(conn)

	// cron service
	s := gocron.NewScheduler(time.UTC)

	reminderSrv := reminderService.NewReminderSrv(s, remindRepo, notificationSrv)

	if firebaseApp != nil {
		reminderSrv.ScheduleNotificationEverySixHours()
		reminderSrv.ScheduleNotificationDaily()
	}

	// reminder service and implementation
	s.Every(5).Minutes().Do(func() {
		log.Println("checking for 5 minutes reminders")
		reminderSrv.SetReminderEvery5Min()
	})

	s.Every(30).Minutes().Do(func() {
		log.Println("checking for 30 minutes reminders")
		reminderSrv.SetReminderEvery30Min()
	})

	// run cron jobs
	s.StartAsync()

	// token service
	srv := tokenservice.NewTokenSrv(secret)

	//logger service
	logger := log_4_go.NewLogger()

	//crypto service
	cryptoSrv := cryptoService.NewCryptoSrv()

	awsSrv := awsService.NewAWSSrv(s3session)

	//email service
	emailSrv := emailService.NewEmailSrv(fromEmailAddr, smtpPWD, smtpHost, smtpPort)

	//Notification Service
	//Note Handle Unable to Connect to Firebase

	//project service
	projectSrv := projectService.NewProjectSrv(projectRepo, timeSrv, validationSrv, logger)

	// task service
	taskSrv := taskService.NewTaskSrv(taskRepo, timeSrv, validationSrv, logger, reminderSrv, notificationSrv)

	// user service

	userSrv := userService.NewUserSrv(userRepo, validationSrv, timeSrv, cryptoSrv, emailSrv, awsSrv, srv, emitter)

	//call service
	callSrv := callService.NewCallSrv(callRepo, timeSrv, validationSrv, logger)

	// social login service

	loginSrv := socialLoginService.NewLoginSrv(userRepo, timeSrv, srv)

	// va service
	vaSrv := vaService.NewVaService(vaRepo, validationSrv, timeSrv, cryptoSrv)

	// subscribe service
	subscribeSrv := subscribeService.NewSubscribeSrv(subRepo, emailSrv, emitter)

	// data service
	dataSrv := dataService.NewDataService(dataRepo)

	r := gin.New()
	r.MaxMultipartMemory = 1 << 20
	r.Use(middlewares.CORS())
	v1 := r.Group("/api/v1")

	// Middlewares
	v1.Use(gin.Logger())
	v1.Use(gin.Recovery())
	v1.Use(gzip.Gzip(gzip.DefaultCompression))

	//handle cors
	//v1.Use(cors.New(cors.Config{
	//	AllowAllOrigins: true,
	//}))

	// routes

	//ping route
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//welcome message route
	v1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Ticked Backend Server - V1.0.0")
	})

	//handle user routes
	routes.UserRoutes(v1, userSrv, srv)

	//handle call routes
	routes.CallRoute(v1, callSrv)

	//handle social login route
	routes.SocialLoginRoute(v1, loginSrv)

	//project routes
	routes.ProjectRoutes(v1, projectSrv, srv)

	//handle task routes
	routes.TaskRoutes(v1, taskSrv, srv)

	//handle Notifications
	routes.NotificationRoutes(v1, notificationSrv, srv)

	//handle VA
	routes.VARoutes(v1, vaSrv, srv, taskSrv, userSrv)

	//handle subscribe route
	routes.SubscribeRoutes(v1, subscribeSrv)

	//handle data route
	routes.DataRoutes(v1, dataSrv)

	// Payment route
	v1.POST("/checkout", paymentHandler.CheckoutCreator)
	v1.POST("/eventService", paymentHandler.HandleEvent)

	//chat service connection
	pusherClient := pusher.Client{
		AppID:   "1512808",
		Key:     "f79030d90753a91854e6",
		Secret:  "06b8abef8713abd21cc9",
		Cluster: "eu",
		Secure:  true,
	}

	v1.POST("dashboard/assistant", func(c *gin.Context) {
		// var data map[string]string
		var data map[string]string

		if err := c.BindJSON(&data); err != nil {
			return
		}
		pusherClient.Trigger("vachat", "message", data)

		c.JSON(http.StatusOK, []string{})
	})

	// Notifications
	// Register to Receive Notifications
	//v1.POST("/notification", notificationHandler.RegisterForNotifications)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	// Documentation
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srvDetails := http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     r,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		log.Println("SERVER STARTING ON PORT:", port)
		err := srvDetails.ListenAndServe()
		if err != nil {
			log.Printf("ERROR STARTING SERVER: %v", err)
			os.Exit(1)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Printf("Closing now, We've gotten signal: %v", sig)

	ctx := context.Background()
	srvDetails.Shutdown(ctx)
}

func connect() (*ampq.Connection, error) {
	var count int64
	var backoff = 1 * time.Second
	var connection *ampq.Connection

	// don't continue until connects

	for {
		c, err := ampq.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			log.Println("Not Connected yet....")
			count++
		} else {
			connection = c
			break
		}

		if count > 5 {
			log.Println(" error connecting", nil)
			return nil, err
		} else {
			backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
			log.Println("waiting....")
			time.Sleep(backoff)
			continue
		}
	}
	log.Println("connected to rabbitMQ")
	return connection, nil
}
