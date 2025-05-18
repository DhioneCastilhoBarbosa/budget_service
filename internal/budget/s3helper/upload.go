package s3helper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func UploadFileToS3(c *gin.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}))

	uploader := s3manager.NewUploader(sess)

	key := "uploads_budgets/" + time.Now().Format("20060102_150405") + "_" + filepath.Base(file.Filename)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   src,
		//ACL:    aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func UploadReaderToS3(reader io.Reader, filename string) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}))

	uploader := s3manager.NewUploader(sess)

	key := "uploads_budgets/" + time.Now().Format("20060102_150405") + "_" + filepath.Base(filename)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
