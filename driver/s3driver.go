package driver

import (
	"crypto/tls"
	"errors"
	"github.com/antigloss/go/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fclairamb/ftpserver/server"
	"io"
	"os"
	"time"
)

// S3 Driver

type S3Driver struct {
	client       *s3.S3
	port         int
	username     string
	password     string
	s3BucketName string
}

func NewS3Driver(
	client *s3.S3,
	port int,
	username, password string,
	s3BucketName string) S3Driver {
	return S3Driver{
		client,
		port,
		username,
		password,
		s3BucketName}
}

func (s *S3Driver) GetSettings() *server.Settings {
	return &server.Settings{
		ListenHost:     "localhost",
		ListenPort:     s.port,
		MaxConnections: 3,
	}
}

func (s *S3Driver) WelcomeUser(cc server.ClientContext) (string, error) {
	return "Welcome to S3 FTP Server", nil
}

func (s *S3Driver) UserLeft(cc server.ClientContext) {
}

func (s *S3Driver) AuthUser(cc server.ClientContext, user, pass string) (server.ClientHandlingDriver, error) {
	if user != s.username || pass != s.password {
		return nil, errors.New("Invalid username or password")
	}
	client := S3ClientDriver{s.client, s.s3BucketName}
	return &client, nil
}

func (s *S3Driver) GetTLSConfig() (*tls.Config, error) {
	return nil, errors.New("Why do you want tls!?")
}

// S3 Client Driver

type S3ClientDriver struct {
	s3Client     *s3.S3
	s3BucketName string
}

func (s *S3ClientDriver) ChangeDirectory(cc server.ClientContext, directory string) error {
	return nil
}

func (s *S3ClientDriver) MakeDirectory(cc server.ClientContext, directory string) error {
	return nil
}

func (s *S3ClientDriver) ListFiles(cc server.ClientContext) ([]os.FileInfo, error) {
	loi := s3.ListObjectsInput{
		Bucket: aws.String(s.s3BucketName),
	}
	loo, err := s.s3Client.ListObjects(&loi)
	if err != nil {
		logger.Warn("ListFiles: %s", err.Error())
		return nil, err
	}
	var files []os.FileInfo
	for _, o := range loo.Contents {
		files = append(files, NewS3FileInfo(*o.Key, *o.Size, *o.LastModified))
		//o.Key
	}
	return files, nil
}

func (s *S3ClientDriver) OpenFile(cc server.ClientContext, path string, flag int) (server.FileStream, error) {
	return &S3File{
		content:    []byte("the data"),
		readOffset: 0,
	}, nil
}

func (s *S3ClientDriver) DeleteFile(cc server.ClientContext, path string) error {
	return errors.New("Forbidden")
}

func (s *S3ClientDriver) GetFileInfo(cc server.ClientContext, path string) (os.FileInfo, error) {
	return nil, errors.New("No file info for you")
}

func (s *S3ClientDriver) RenameFile(cc server.ClientContext, from, to string) error {
	return errors.New("Rename not allowed. Not for you")
}

func (s *S3ClientDriver) CanAllocate(cc server.ClientContext, size int) (bool, error) {
	return false, errors.New("Allocate this.")
}

func (s *S3ClientDriver) ChmodFile(cc server.ClientContext, path string, mode os.FileMode) error {
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
	content    []byte
	readOffset int
}

func (f *S3File) Close() error {
	return nil
}

func (f *S3File) Read(buffer []byte) (int, error) {
	n := copy(buffer, f.content[f.readOffset:])
	f.readOffset += n
	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

func (f *S3File) Seek(n int64, w int) (int64, error) {
	return 0, nil
}

func (f *S3File) Write(buffer []byte) (int, error) {
	return 0, nil
}
