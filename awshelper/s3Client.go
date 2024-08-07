package awshelper

import (
	"context"
	"encoding/base64"
	"log"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const(
	BUCKET_NAME = "abhis-s3"
	BUCKET_REGION = "ap-south-1"
)

var (
	S3ServiceObject *S3Service
)

func init(){

	accessKey := "add accessKey"
	secretKey := "add secret"
	options := s3.Options{
		Region:     BUCKET_REGION,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		ClientLogMode: aws.LogRetries|aws.LogResponseWithBody,
		RetryMaxAttempts: 1,
	}

	S3ServiceObject = &S3Service{
		S3Client: s3.New(options),
	}
}


type S3Service struct {
	S3Client *s3.Client
}

func (service S3Service) UploadBase64(b64 string, key string)  error{

	data := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	
	size := calcOrigBinaryLength(b64)
	log.Println("file size is : ", size)

	_, err := service.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(key),
		Body: data,
		ContentLength: aws.Int64(int64(size)),
		ContentType: aws.String(".png"),
	})

	if  err != nil {
		log.Println("uploadBase64 : error in thumbnail upload : " , err)
		return err
	}

	return nil
}

func calcOrigBinaryLength(datas string) int {
    l := len(datas)

    eq := 0
    if l >= 2 {
        if datas[l-1] == '=' {
            eq++
        }
        if datas[l-2] == '=' {
            eq++
        }

        l -= eq
    }

    return (l*3 - eq) / 4
}



