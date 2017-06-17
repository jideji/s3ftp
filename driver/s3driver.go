package driver

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fclairamb/ftpserver/server"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// S3 Driver

type S3Driver struct {
	client       *s3.S3
	host         string
	port         int
	username     string
	password     string
	s3BucketName string
}

func NewS3Driver(
	client *s3.S3,
	host string,
	port int,
	username, password string,
	s3BucketName string) S3Driver {
	return S3Driver{
		client,
		host,
		port,
		username,
		password,
		s3BucketName}
}

func (s *S3Driver) GetSettings() *server.Settings {
	return &server.Settings{
		ListenHost:     s.host,
		ListenPort:     s.port,
		MaxConnections: 100,
	}
}

func (s *S3Driver) WelcomeUser(cc server.ClientContext) (string, error) {
	log.Debug("Welcome message")
	return "Welcome to S3 FTP Server", nil
}

func (s *S3Driver) UserLeft(cc server.ClientContext) {
	log.Info("User left")
}

func (s *S3Driver) AuthUser(cc server.ClientContext, user, pass string) (server.ClientHandlingDriver, error) {
	if user != s.username || pass != s.password {
		log.Warn("Invalid username or password")
		return nil, errors.New("Invalid username or password")
	}
	log.Info("Valid username and password")
	client := S3ClientDriver{s.client, s.s3BucketName}
	return &client, nil
}

func (s *S3Driver) GetTLSConfig() (*tls.Config, error) {
	log.Warn("TLS config")
	return nil, errors.New("Why do you want tls!?")
}

// S3 Client Driver

type S3ClientDriver struct {
	s3Client     *s3.S3
	s3BucketName string
}

func (s *S3ClientDriver) ChangeDirectory(cc server.ClientContext, directory string) error {
	log.Info("Change directory (ignored)")
	return nil
}

func (s *S3ClientDriver) MakeDirectory(cc server.ClientContext, directory string) error {
	log.Info("Make directory (ignored)")
	return nil
}

func (s *S3ClientDriver) ListFiles(cc server.ClientContext) ([]os.FileInfo, error) {
	log.Info("List files")
	loi := s3.ListObjectsInput{
		Bucket: aws.String(s.s3BucketName),
	}
	loo, err := s.s3Client.ListObjects(&loi)
	if err != nil {
		log.Warn("ListFiles: %s", err.Error())
		return nil, err
	}
	var files []os.FileInfo
	for _, o := range loo.Contents {
		files = append(files, NewS3FileInfo(*o.Key, *o.Size, *o.LastModified))
	}
	return files, nil
}

func (s *S3ClientDriver) OpenFile(cc server.ClientContext, path string, flag int) (server.FileStream, error) {
	if flag != os.O_WRONLY {
		log.Warn("Open file with wrong flag %d", flag)
		return nil, errors.New("Only writing files supported")
	}

	log.Info("Open file with path %s", path)

	return &S3File{
		s3Client:   s.s3Client,
		bucketName: s.s3BucketName,
		path:       path,
		content:    []byte{},
		readOffset: 0,
	}, nil
}

func (s *S3ClientDriver) DeleteFile(cc server.ClientContext, path string) error {
	log.Warn("Delete file %s (no-op)", path)
	return nil
}

func (s *S3ClientDriver) GetFileInfo(cc server.ClientContext, path string) (os.FileInfo, error) {
	log.Warn("Get file info %s", path)
	return nil, errors.New("No file info for you")
}

func (s *S3ClientDriver) RenameFile(cc server.ClientContext, from, to string) error {
	log.Warn("Rename file from %s to %s", from, to)
	return errors.New("Rename not allowed. Not for you")
}

func (s *S3ClientDriver) CanAllocate(cc server.ClientContext, size int) (bool, error) {
	log.Warn("Can allocate %d", size)
	return false, errors.New("Allocate this.")
}

func (s *S3ClientDriver) ChmodFile(cc server.ClientContext, path string, mode os.FileMode) error {
	log.Warn("Chmod %s %d", path, mode)
	return errors.New("Forbidden")
}

// S3 FileInfo

func NewS3FileInfo(name string, size int64, lastModified time.Time) *S3FileInfo {
	return &S3FileInfo{
		name:         name,
		size:         size,
		lastModified: lastModified,
	}
}

type S3FileInfo struct {
	name         string
	size         int64
	lastModified time.Time
}

func (s *S3FileInfo) Name() string {
	return s.name
}

func (s *S3FileInfo) Size() int64 {
	return s.size
}

func (s *S3FileInfo) Mode() os.FileMode {
	return 0444
}

func (s *S3FileInfo) ModTime() time.Time {
	return s.lastModified
}

func (s *S3FileInfo) IsDir() bool {
	return false
}

func (s *S3FileInfo) Sys() interface{} {
	return nil
}

// S3 File

type S3File struct {
	s3Client   *s3.S3
	bucketName string
	path       string
	content    []byte
	readOffset int
}

func (f *S3File) Close() error {
	log.Info("Sending file s3://%s%s to s3", f.bucketName, f.path)
	poi := s3.PutObjectInput{
		Bucket: aws.String(f.bucketName),
		Key:    aws.String(f.path),
		Body:   bytes.NewReader(f.content),
	}
	if _, err := f.s3Client.PutObject(&poi); err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (f *S3File) Read(buffer []byte) (int, error) {
	return 0, nil
}

func (f *S3File) Seek(n int64, w int) (int64, error) {
	return 0, nil
}

func (f *S3File) Write(buffer []byte) (int, error) {
	f.content = append(f.content, buffer...)
	return len(buffer), nil
}
