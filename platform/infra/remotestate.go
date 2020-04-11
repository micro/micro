package infra

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// RemoteState is an action that verifies remote state is available
type RemoteState struct {
	ID   string
	Name string
}

// Validate checks the remote state buckets and table exist
func (r *RemoteState) Validate() error {
	if err := r.validateConfig(); err != nil {
		fmt.Fprintf(os.Stdout, "[%s] The remote state backend is invalid!\n", r.Name)
		return err
	}
	fmt.Fprintf(os.Stdout, "[%s] The remote state backend is valid\n", r.Name)
	return nil
}

// Plan does nothing
func (r *RemoteState) Plan() error {
	return nil
}

// Apply does nothing
func (r *RemoteState) Apply() error {
	return nil
}

// Finalise does nothing
func (r *RemoteState) Finalise() error {
	return nil
}

// Destroy does nothing
func (r *RemoteState) Destroy() error {
	return nil
}

func (r *RemoteState) validateConfig() error {
	stateStore := viper.GetString("state-store")
	switch stateStore {
	case "aws":
		return r.validateAws()
	case "azure":
		return r.validateAzure()
	default:
		return errors.Errorf("%s is not a supported state store", stateStore)
	}
}

func (r *RemoteState) validateAzure() error {
	//TODO: meaningful validation
	return nil
}

func (r *RemoteState) validateAws() error {
	client := s3.New(
		session.New(
			&aws.Config{
				Region: func() *string {
					if r := os.Getenv("AWS_REGION"); len(r) != 0 {
						return &r
					}
					return aws.String("eu-west-2")
				}(),
			},
		),
	)
	if _, err := client.PutObject(
		&s3.PutObjectInput{
			Key:    aws.String(r.ID),
			Bucket: aws.String(viper.GetString("aws-s3-bucket")),
			Body:   strings.NewReader(r.ID),
		},
	); err != nil {
		return errors.Wrap(err, "Could not put an object in to the remote state bucket")
	}

	read, err := client.GetObject(
		&s3.GetObjectInput{
			Key:    aws.String(r.ID),
			Bucket: aws.String(viper.GetString("aws-s3-bucket")),
		},
	)
	if err != nil {
		return errors.Wrap(err, "Could not read back an object from to the remote state bucket")
	}
	body, err := ioutil.ReadAll(read.Body)
	read.Body.Close()
	if err != nil {
		return errors.Wrap(err, "Error reading body from S3")
	}
	if string(body) != r.ID {
		return fmt.Errorf("Read back an invalid value from remote state. Expected %s, got %s", r.ID, string(body))
	}

	if _, err := client.DeleteObject(
		&s3.DeleteObjectInput{
			Key:    aws.String(r.ID),
			Bucket: aws.String(viper.GetString("aws-s3-bucket")),
		},
	); err != nil {
		return errors.Wrap(err, "Error deleting object from S3")
	}
	return nil
}
