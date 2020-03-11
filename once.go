package once

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

func (m *manager) Ensure(req events.APIGatewayProxyRequest) error {
	requestId := req.RequestContext.RequestID

	if requestId == ""{
		return fmt.Errorf("api gateway request id not exist")
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(m.s3Bucket),
		Key: aws.String(path.Join(m.s3Prefix, requestId)),
		ObjectLockLegalHoldStatus:aws.String(s3.ObjectLockLegalHoldStatusOn),
	}

	s3Service := s3.New(session.Must(session.NewSession()))

	_, err := s3Service.PutObject(putInput)

	return err
}
