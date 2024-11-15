package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appConfig "github.com/toufiq-austcse/deployit/config"
	"os"
)

type S3ManagerService struct {
	client *s3.Client
}

func NewS3ManagerService() (*S3ManagerService, error) {
	creds := credentials.NewStaticCredentialsProvider(appConfig.AppConfig.AWS_CONFIG.ACCESS_KEY_ID, appConfig.AppConfig.AWS_CONFIG.SECRET_ACCESS_KEY, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appConfig.AppConfig.AWS_CONFIG.REGION),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	fmt.Println("S3ManagerService Initialized")
	return &S3ManagerService{client: client}, nil
}

func (s3ManagerService S3ManagerService) UploadFile(filePath string, s3Folder string) (*string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s", s3Folder, fileInfo.Name())
	_, putErr := s3ManagerService.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &appConfig.AppConfig.AWS_CONFIG.BUCKET_NAME,
		Key:    &key,
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if putErr != nil {
		return nil, err
	}

	if removeErr := os.Remove(filePath); removeErr != nil {
		fmt.Println("error in removing file ", removeErr.Error())
	}
	return &key, nil

}
