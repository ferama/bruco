package getter

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Getter adds support for function loading from minio s3
// It enables stuff like this:
// $ export AWS_ACCESS_KEY_ID=minioadmin
// $ export AWS_SECRET_ACCESS_KEY=minioadmins
// $ bruco s3://localhost:9000/bruco/sentiment.zip
type S3Getter struct {
	getterCommon
}

// NewS3Getter builds a new s3getter object
func NewS3Getter() *S3Getter {
	return &S3Getter{}
}

// Download
func (g *S3Getter) Download(resourceURL string) (string, error) {
	file, err := g.getTmpFile(resourceURL)
	if err != nil {
		return "", fmt.Errorf("can't store file %s", err)
	}
	defer file.Close()
	g.PayloadPath = file.Name()

	secure := true
	parsed, _ := url.Parse(resourceURL)
	if strings.ToLower(parsed.Scheme) == "s3" {
		secure = false
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := parsed.Host

	parts := strings.Split(parsed.Path, "/")
	bucket := parts[1]
	objectKey := strings.Join(parts[2:], "/")

	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return "", err
	}
	if err := s3Client.FGetObject(context.Background(),
		bucket, objectKey, file.Name(), minio.GetObjectOptions{}); err != nil {
		return "", err
	}

	return file.Name(), nil
}
