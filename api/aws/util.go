package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
		log.Warn(err)
		panic(err)
	}
	return ses
}
