package utils

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/gommon/log"
)

func InitS3(region, id, secret string) *session.Session {
	var ses, err = session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				id, secret, "",
			),
		},
	)

	if err != nil {
		log.Error(err)
	}
	return ses
}

func UploadFileToS3(ses *session.Session, fileHeader multipart.FileHeader) string {

	var manager = s3manager.NewUploader(ses)
	var src, err = fileHeader.Open()
	if err != nil {
		log.Info(err)
	}
	defer src.Close()

	size := fileHeader.Size
	buffer := make([]byte, size)
	src.Read(buffer)

	var res, err1 = manager.Upload(
		&s3manager.UploadInput{
			Bucket:       aws.String("karen-givi-bucket"),
			Key:          aws.String(fileHeader.Filename),
			ACL:          aws.String("public-read-write"),
			Body:         bytes.NewReader(buffer),
			ContentType:  aws.String(http.DetectContentType(buffer)),
			StorageClass: aws.String("STANDARD"),
		},
	)

	if err1 != nil {
		log.Info(res)
		log.Error(err1)
	}

	var url = "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/" + fileHeader.Filename

	return url
}
