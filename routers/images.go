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
	"github.com/starlingapps/twittergo/awsgo"
	"github.com/starlingapps/twittergo/db"
	"github.com/starlingapps/twittergo/models"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadImage(req events.APIGatewayProxyRequest, userId string) models.ApiResponse {
	var res = models.ApiResponse{
		Status: 400,
	}
	fmt.Printf("> Into UploadImage method for userId %s\n", userId)
	user, err := db.FindUserById(userId)
	if err != nil {
		res.Status = 500
		res.Body = "UserProfile for user not found"
		return res
	}
	bucketName := os.Getenv("BUCKET_NAME")
	imageType := req.QueryStringParameters["imageType"]
	var objectKey string
	switch imageType {
	case "Avatar":
		objectKey = user.Avatar
	case "Banner":
		objectKey = user.Banner
	default:
		res.Body = "Image type unknown, it should be: Avatar or Banner"
		return res
	}
	client := s3.NewFromConfig(awsgo.Cfg)
	sClient := s3.NewPresignClient(client)
	presignUrl, err := sClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(time.Minute*15))

	res.Status = 200
	res.Body = presignUrl.URL

	return res

	// The above process get a signed URL to avoid more processing on the lambda.
	// The below commented lines are to return the image file content to client
	//  it is easier to return a public URL and visualize the image in Postman

	/*
			result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			})
			if err != nil {
				res.Status = 500
				res.Body = fmt.Sprintf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
				fmt.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
				return res
			}
			defer result.Body.Close()
			bytesFile, err := io.ReadAll(result.Body)
			if err != nil {
				res.Status = 500
				res.Body = fmt.Sprintf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
				fmt.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
				return res
			}
			buffer := bytes.NewBuffer(bytesFile)
			res.ActualResponse = &events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       buffer.String(),
				Headers: map[string]string{
					"Content-Type":        "application/octet-stream",
					"Content-Disposition": fmt.Sprintf("attachment: filename=\"%s\"", objectKey),
				},
			}
			res.Status = 200
			res.Body = "Image found"
		return res
	*/
}

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
			ext := filepath.Ext(imagePart.FileName())
			name := strings.TrimSuffix(objectKey, filepath.Ext(objectKey))
			objectKey = name + ext
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
		user := models.User{}
		user.Avatar = objectKey
		status, err := db.UpdateProfileById(userId, user)
		if err != nil || !status {
			res.Status = 500
			res.Body = "Error updating Avatar for user"
			fmt.Printf("Error updating Avatar for user %s\n", err)
			return res
		}

		res.Status = 200
		res.Body = fmt.Sprintf("%s updated correctly", imageType)
		return res

	} else {
		res.Status = 400
		res.Body = "Image must sent into multipart-form-data form"
		return res
	}
}
