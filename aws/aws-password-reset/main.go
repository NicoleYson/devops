package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type iamClient struct {
	svc iamiface.IAMAPI
}

func main() {
	aws := iamClient{}
	if aws.svc == nil {
		aws.svc = aws.newIamClient()
	}
	aws.resetPassword()
}

// newIamClient will return a IAM client.
func (client *iamClient) newIamClient() iamiface.IAMAPI {
	sess := session.Must(session.NewSession())
	if *sess.Config.Region == "" {
		sess.Config.Region = aws.String("us-west-2")
	}

	return iam.New(sess, &aws.Config{})
}
