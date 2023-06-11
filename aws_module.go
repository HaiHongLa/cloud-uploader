package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func aws_config() bool{
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

func aws_upload(region string, bucketName string, filePath string, fileKey string) bool {
	configFile, err := os.Open("cloudConfig/AWSConfig")
	if err != nil {
		fmt.Println("AWS Configuration not found\nConfigure your AWS S3 credentials by running: file-storage configure aws", err)
		return false
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
		fmt.Println("AWS access key is missing. Please run: file-storage configure aws")
		return false
	}
	if keys["AWS_SECRET_ACCESS_KEY"] == "" {
		fmt.Println("Secret access key is missing. Please run: file-storage configure aws")
		return false
	}

	accessKey := keys["AWS_ACCESS_KEY"]
	secretAccessKey := keys["AWS_SECRET_ACCESS_KEY"]

	// Upload the file
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey, 
			secretAccessKey, 
			""),
	}))
	svc := s3.New(sess)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file", err)
		return false
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileKey),
		Body: file,
	})
	if err != nil {
		fmt.Println("Failed to upload file to S3", err)
		return false
	}

	fmt.Println("File uploaded successfully")
	return true
}