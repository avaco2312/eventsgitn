package aws

import (
	awsp "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var Sesion *session.Session

func SetSession() error {
	var err error
	if Sesion != nil {
		return nil
	}
	Sesion, err = session.NewSession(&awsp.Config{
		Region: awsp.String("us-west-2"),
	})
	return err
}
