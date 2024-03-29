package awsService

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSService interface {
	UploadImage(file multipart.File, filename string) error
}

type awsSrv struct {
	s3session *s3.S3
}

func (a *awsSrv) UploadImage(file multipart.File, filename string) error {

	_, err := a.s3session.PutObject(&s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String("ticked-v1-backend-bucket"),
		Key:         aws.String(filename),
		ContentType: aws.String("image/jpeg"),
		// ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})

	if err != nil {
		return err
	}

	return nil
}

func NewAWSSrv(s *s3.S3) AWSService {
	return &awsSrv{s3session: s}
}
