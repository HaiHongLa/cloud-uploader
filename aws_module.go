package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func configureAWSS3() bool {
	fmt.Println("Enter Access Key: ")
	var accessKey string
	fmt.Scanln(&accessKey)
	if len(accessKey) == 0 {
		fmt.Println("Please enter your access key")
		return false
	}
	fmt.Println("Enter Secret Access Key: ")
	var secretAccessKey string
	fmt.Scanln(&secretAccessKey)
	if len(secretAccessKey) == 0 {
		fmt.Println("Please enter your secret access key")
		return false
	}
	err := os.MkdirAll("cloudConfig", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return false
	}
	awsConfig := "AWS_ACCESS_KEY=" + accessKey + "\n" + "AWS_SECRET_ACCESS_KEY=" + secretAccessKey
	err = ioutil.WriteFile("cloudConfig/AWSConfig", []byte(awsConfig), 0644)
	if err != nil {
		fmt.Println("Error creating file", err)
		return false
	}
	fmt.Println("Configured AWS S3 successfully")
	return true
}

func getAWSS3Session(region string) (*session.Session, error) {
	configFile, err := os.Open("cloudConfig/AWSConfig")
	if err != nil {
		return nil, fmt.Errorf("AWS configuration not found. Configure your AWS S3 credentials by running: file-storage configure aws. Error: %v", err)
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	keys := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.Split(line, "=")
		if len(pair) != 2 {
			continue
		}
		keys[pair[0]] = pair[1]
	}
	if keys["AWS_ACCESS_KEY"] == "" {
		return nil, fmt.Errorf("AWS access key is missing. Please run: file-storage configure aws")
	}
	if keys["AWS_SECRET_ACCESS_KEY"] == "" {
		return nil, fmt.Errorf("Secret access key is missing. Please run: file-storage configure aws")
	}

	accessKey := keys["AWS_ACCESS_KEY"]
	secretAccessKey := keys["AWS_SECRET_ACCESS_KEY"]

	// Create an AWS session with the provided credentials
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretAccessKey,
			""),
	}))
	return sess, nil
}

func getAWSS3Client(region string) (*s3.S3, error) {
	sess, err := getAWSS3Session(region)
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS S3 session: %v", err)
	}
	svc := s3.New(sess)
	return svc, nil
}

func uploadToAWSS3Bucket(region string, bucketName string, filePath string, fileKey string) bool {
	svc, err := getAWSS3Client(region)
	if err != nil {
		fmt.Println("Failed to authenticate: ", err)
		return false
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open local file", err)
		return false
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   file,
	})
	if err != nil {
		fmt.Println("Failed to upload file to S3", err)
		return false
	}

	fmt.Println("File uploaded successfully")
	return true
}

func listAllFilesFromAWSS3Bucket(region string, bucketName string) bool {
	svc, err := getAWSS3Client(region)
	if err != nil {
		fmt.Println("Failed to authenticate: ", err)
		return false
	}
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("Failed to list objects: ", err)
		return false
	}
	for _, obj := range result.Contents {
		fmt.Println(*obj.Key)
	}
	return true
}

func deleteFromS3Bucket(region string, bucketName string, fileKey string) bool {
	svc, err := getAWSS3Client(region)
	if err != nil {
		fmt.Println("Failed to authenticate: ", err)
		return false
	}
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	}

	// Delete the file from S3
	_, err = svc.DeleteObject(input)
	if err != nil {
		fmt.Println("Failed to delete file from S3: ", err)
		return false
	}
	fmt.Println("Deleted file " + fileKey + " successfully")
	return true
}

func getFileFromS3(region string, bucketName string, fileKey string, filePath string) bool {
	sess, err := getAWSS3Session(region)
	if err != nil {
		fmt.Printf("Error creating AWS S3 session: %v\n", err)
		return false
	}
	downloader := s3manager.NewDownloader(sess)
	dirPath := filepath.Dir(filePath)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return false
	}
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return false
	}
	defer file.Close()

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		fmt.Println("Failed to download file:", err)
		return false
	}

	fmt.Println("File downloaded successfully!")
	return true
}
