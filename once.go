package once

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"path"
)

type manager struct {
	s3Bucket string
	s3Prefix string

	// TODO auth
}

func New(bucket string, ops ...option) *manager {
	result := &manager{
		s3Bucket: bucket,
	}

	for _, o := range ops {
		o(result)
	}

	return result
}

func Prefix(p string) option {
	return func(manager *manager) {
		manager.s3Prefix = p
	}
}

type option func(m *manager)

func (m *manager) Ensure(id string) {

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(m.s3Bucket),
		Key: aws.String(path.Join(m.s3Prefix, id)),
		ObjectLockLegalHoldStatus:aws.String(s3.ObjectLockLegalHoldStatusOn),
	}

	s3Service := s3.New(session.Must(session.NewSession()))

	_, err := s3Service.PutObject(putInput)

	if err != nil{
		panic(err)
	}
}
