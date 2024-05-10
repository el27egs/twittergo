package routers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	m "github.com/grokify/go-awslambda"
	"github.com/starlingapps/twittergo/models"
	"io"
	"mime"
	"os"
	"strings"
)

func UploadImage(req events.APIGatewayProxyRequest, userId string) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into UploadImage method for userId %s\n", userId)
	bucketName := os.Getenv("BUCKET_NAME")
	imageType := req.QueryStringParameters["imageType"]
	var objectKey string
	switch imageType {
	case "Avatar":
		objectKey = "avatars/" + userId + ".jpg"
	case "Banner":
		objectKey = "banners/" + userId + ".jpg"
	default:
		res.Body = "Image type unknown, it should be: Avatar or Banner"
		return res
	}
	// Validate that Content-Type value is multipart/form-data
	mediaType, _, err := mime.ParseMediaType(req.Headers["Content-Type"])
	if err != nil {
		res.Status = 500
		res.Body = "Form data type unknown, is must be only multipart/form-data"
		fmt.Printf("Form data type unknown, is must be only multipart/form-data %s", err)
		return res
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		// Reading part using boundaries to get the part corresponding to image data
		// This library has almost same lines of code used by the trainer, so
		// I can say it follows same logic as the videos.
		bodyReader, err := m.NewReaderMultipart(req)
		imagePart, err := bodyReader.NextPart()
		if err != nil && err != io.EOF {
			res.Status = 500
			res.Body = "Error reading image data"
			fmt.Printf("Form reading image data %s\n", err)
			return res
		}
		// Creating a Reader using buffer to create the file in memory to pass to s3
		if imagePart.FileName() != "" {
			byteArrayImage, err := io.ReadAll(imagePart)
			imageReader := bytes.NewBuffer(byteArrayImage)

			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
			if err != nil {
				res.Status = 500
				res.Body = fmt.Sprintf("unable to load SDK config, %s", err)
				fmt.Printf("unable to load SDK config, %s", err)
				return res
			}
			client := s3.NewFromConfig(cfg)
			_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   imageReader,
			})
			if err != nil {
				res.Status = 500
				res.Body = fmt.Sprintf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					imageReader, bucketName, objectKey, err)
				fmt.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					imageReader, bucketName, objectKey, err)
				return res
			}
		}
		res.Status = 200
		res.Body = "Image uploaded!!!"
		return res

	} else {
		res.Status = 200
		res.Body = "Image must sent into multipart-form-data form"
		return res
	}
}
