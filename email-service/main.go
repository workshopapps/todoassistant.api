package main

import (
	mail_grpc "emailSrv/cmd/mail-grpc"
	mail_grpc_mail "emailSrv/internal/grpc-mail"
	mail_serivce "emailSrv/internal/service/mail-serivce"
	utility "emailSrv/internal/utils"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

func main() {
	config, err := utility.LoadConfig("./")
	if err != nil {
		return
	}
	// get Data from env file
	password := config.Password
	if password == "" {
		log.Fatalf("Unable to retrieve password from env file")
	}

	host := config.Host
	if host == "" {
		log.Fatalf("Unable to retrieve host from env file")
	}

	port := config.MailPort
	if port == "" {
		log.Fatalf("Unable to retrieve port from env file")
	}

	fromEmail := config.FromEmail
	if fromEmail == "" {
		log.Fatalf("Unable to retrieve sender Email from env file")
	}
	grpcPort := config.GrpcPort
	if grpcPort == "" {
		log.Fatalf("Unable to retrieve sender Email from env file")
	}

	// create new Mailing Service
	srv := mail_serivce.NewEmailSrv(fromEmail, password, host, port)

	//start up grpc server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for grpc: %v", err)
	}
	s := grpc.NewServer()
	mail_grpc_mail.RegisterMailServiceServer(s, mail_grpc.NewMailGrpcServer(srv))
	log.Println("grpc server starting...")

	// listen for grpc connections
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to listen for grpc: %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)   // create a channel that can for signals
	signal.Notify(sigChan, os.Interrupt) // whenever an interrupt signal is received send to the channel
	signal.Notify(sigChan, os.Kill)      // whenever a kill signal is received send to the channel

	sig := <-sigChan // read from the channel (this would block until there is something to read in the channel)
	log.Printf("Closing now, We've gotten signal: %v", sig)
	// graceful shutdown
	s.GracefulStop()
}
