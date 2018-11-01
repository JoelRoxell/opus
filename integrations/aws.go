package integrations

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AwsService is used to access different services found in aws.
type AwsService struct {
	session  *session.Session
	uploader *s3manager.Uploader
	Region   string
}

// Init constructs a AwsService with a new AWS session.
func (a *AwsService) Init() {
	session, _ := session.NewSession(&aws.Config{Region: &a.Region})

	a.session = session
	a.uploader = s3manager.NewUploader(session)
}

// UploadFile uploads a file given by the given filepath, to a specific bucket.
func (a *AwsService) UploadFile(filePath string, bucket string) error {
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

	body, err := ioutil.ReadAll(filePointer)
	byteReader := bytes.NewReader(body)

	config := &s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &filePath,
		Body:        byteReader,
		ContentType: &contentType,
	}

	res, err := a.uploader.Upload(config)

	if err != nil {
		return err
	}

	log.Printf("Uploaded %s to %s; %s ContentType: %s", filePath, res.Location, res.UploadID, contentType)

	return nil
}
