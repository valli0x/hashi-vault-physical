package s3

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestDefaultS3Backend(t *testing.T) {
	DoS3BackendTest(t, "")
}

func TestS3BackendSseKms(t *testing.T) {
	DoS3BackendTest(t, "alias/aws/s3")
}

func DoS3BackendTest(t *testing.T, kmsKeyId string) {
	if enabled := os.Getenv("VAULT_ACC"); enabled == "" {
		t.Skip()
	}

	logger := logging.NewVaultLogger(log.Debug)

	credsConfig := &awsutil.CredentialsConfig{Logger: logger}

	credsChain, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		t.Fatal(err)
	}

	_, err = credsChain.Get()
	if err != nil {
		t.Fatal(err)
	}

	// If the variable is empty or doesn't exist, the default
	// AWS endpoints will be used
	endpoint := os.Getenv("AWS_S3_ENDPOINT")

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credsChain,
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	})
	if err != nil {
		t.Fatal(err)
	}
	s3conn := s3.New(sess)

	randInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-s3-testacc-%d", randInt)

	_, err = s3conn.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		// Gotta list all the objects and delete them
		// before being able to delete the bucket
		listResp, _ := s3conn.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
		})

		objects := &s3.Delete{}
		for _, key := range listResp.Contents {
			oi := &s3.ObjectIdentifier{Key: key.Key}
			objects.Objects = append(objects.Objects, oi)
		}

		s3conn.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: objects,
		})

		_, err := s3conn.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(bucket)})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// This uses the same logic to find the AWS credentials as we did at the beginning of the test
	b, err := NewS3Backend(map[string]string{
		"bucket":   bucket,
		"kmsKeyId": kmsKeyId,
		"path":     "test/vault",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
