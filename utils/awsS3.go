package utils

import (
	"bytes"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/gommon/log"
	"github.com/lithammer/shortuuid"
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

	var uid = shortuuid.New()

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
			Key:          aws.String(uid),
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

	var url = "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/" + uid

	return url
}

func DeleteFileS3(ses *session.Session, name string) string {
	var svc = s3.New(ses)

	var input = &s3.DeleteObjectInput{
		Bucket: aws.String("karen-givi-bucket"),
		Key:    aws.String(name),
	}

	var res, err = svc.DeleteObject(input)

	if err != nil {
		log.Info(res)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Info(aerr.Error())
			}
			return err.Error()
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Info(err.Error())
			return err.Error()
		}
	}

	return "success"
}

func UpdateFileS3(ses *session.Session, name string, fileHeader multipart.FileHeader) string {

	src, err := fileHeader.Open()
	if err != nil {
		log.Info(err)
		return err.Error()
	}
	defer src.Close()

	size := fileHeader.Size
	buffer := make([]byte, size)
	src.Read(buffer)

	var svc = s3.New(ses)

	var input = &s3.PutObjectInput{
		Body:         bytes.NewReader(buffer),
		Bucket:       aws.String("karen-givi-bucket"),
		Key:          aws.String(name),
		ACL:          aws.String("public-read-write"),
		ContentType:  aws.String(http.DetectContentType(buffer)),
		StorageClass: aws.String("STANDARD"),
	}

	res, err := svc.PutObject(input)

	if err != nil {
		log.Info(res)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Info(aerr.Error())
			}
			return err.Error()
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Info(err.Error())
			return err.Error()
		}
	}

	return "success"
}
