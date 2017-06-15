package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fclairamb/ftpserver/server"
	"github.com/jideji/s3ftp/driver"
	"os"
	"strconv"
)

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Errorf("No environment variable %s defined", key))
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		nValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return nValue
	}
	return defaultValue
}

func main() {
	port := getEnvInt("FTP_PORT", 21)
	s3BucketName := mustGetEnv("S3_BUCKET_NAME")
	ftpUsername := mustGetEnv("FTP_USER")
	ftpPassword := mustGetEnv("FTP_PASS")

	sess := session.Must(session.NewSession())
	creds := credentials.NewEnvCredentials()
	svc := s3.New(sess, &aws.Config{Credentials: creds})
	s3 := driver.NewS3Driver(svc, port, ftpUsername, ftpPassword, s3BucketName)
	ftpServer := server.NewFtpServer(&s3)
	panic(ftpServer.ListenAndServe())
}
