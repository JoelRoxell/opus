package integrations

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AwsService is used to access different services found in aws.
type AwsService struct {
	Session  *session.Session
	Uploader *s3manager.Uploader
	Region   string
}

// Init constructs a AwsService with a new AWS session.
func (a *AwsService) Init() {
	session, _ := session.NewSession(&aws.Config{Region: &a.Region})

	a.Session = session
	a.Uploader = s3manager.NewUploader(session)
}

// UploadFile uploads a file given by the given filepath, to a specific bucket.
func (a *AwsService) UploadFile(filePath string, targetPath string, bucket string) error {
	filePointer, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer filePointer.Close()

	buffer := make([]byte, 512)

	n, err := filePointer.Read(buffer)

	if err != nil {
		fmt.Println("Error:", err)
	}

	contentType := http.DetectContentType(buffer[:n])

	filePointer.Seek(0, 0)
	body, err := ioutil.ReadAll(filePointer)
	byteReader := bytes.NewReader(body)

	config := &s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &targetPath,
		Body:        byteReader,
		ContentType: &contentType,
	}

	res, err := a.Uploader.Upload(config)

	if err != nil {
		return err
	}

	log.Printf("Uploaded %s to %s; %s ContentType: %s", filePath, res.Location, res.UploadID, contentType)

	return nil
}

// Download fetches the specified remote file from bucket and writes it to the a target path.
func (a *AwsService) Download(remotePath string, toPath string, bucket string) {
	downloader := s3manager.NewDownloader(a.Session)

	path, err := filepath.Abs(toPath)

	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		panic(err)
	}

	fileW, err := os.Create(path)

	defer fileW.Close()

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	_, err = downloader.Download(fileW, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &remotePath,
	})

	if err != nil {
		panic(err)
	}
}
