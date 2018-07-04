package actions

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func ImageUpload(c *gin.Context) {
	inputFile, _ := c.FormFile("file")

	creds := credentials.NewStaticCredentials(os.Getenv("ACCESS_TOKEN"), os.Getenv("SECRET"), "")
	_, err := creds.Get()
	if err != nil {
		// handle error
	}
	cfg := aws.NewConfig().WithRegion("ap-south-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	file, err := inputFile.Open()
	if err != nil {
		// handle error
	}
	defer file.Close()

	size := inputFile.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/media/" + inputFile.Filename
	params := &s3.PutObjectInput{
		Bucket:        aws.String("golangdemo"),
		Key:           aws.String(path),
		ACL:           aws.String("public-read"),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		// handle error
	}
	c.JSON(200, gin.H{
		"imagename": inputFile.Filename,
		"url":       awsutil.StringValue(resp)})
}
