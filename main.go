package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fclairamb/ftpserver/server"
	"github.com/jideji/s3ftp/config"
	"github.com/jideji/s3ftp/driver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/inconshreveable/log15.v2"
	"os"
)

func main() {
	log15.Root().SetHandler(log15.StreamHandler(os.Stdout, log15.JsonFormat()))
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	cfg := config.LoadConfig()

	sess := session.Must(session.NewSession())
	creds := credentials.NewEnvCredentials()
	svc := s3.New(sess, &aws.Config{Credentials: creds})
	s3 := driver.NewS3Driver(
		svc,
		cfg.Ftp.Host,
		cfg.Ftp.Port,
		cfg.Ftp.Username,
		cfg.Ftp.Password,
		cfg.S3.BucketName)
	ftpServer := server.NewFtpServer(&s3)
	panic(ftpServer.ListenAndServe())
}
