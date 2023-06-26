package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func configureGCS() bool {
	fmt.Println("Enter path to key file: ")
	var fileKeyPath string
	fmt.Scanln(&fileKeyPath)
	if len(fileKeyPath) == 0 {
		fmt.Println("Please enter path to key file")
		return false
	}
	err := os.MkdirAll("cloudConfig", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return false
	}
	gcsConfig := "GCS_KEY_FILE_PATH=" + fileKeyPath
	err = ioutil.WriteFile("cloudConfig/GCSConfig", []byte(gcsConfig), 0644)
	if err != nil {
		fmt.Println("Error creating file", err)
		return false
	}
	fmt.Println("Configured AWS S3 successfully")
	return true
}

func getGCSClient(ctx *context.Context) (*storage.Client, error) {
	configFile, err := os.Open("cloudConfig/GCSConfig")
	if err != nil {
		fmt.Println("GCS Configuration not found. Configure your Google Cloud Storage credentials by running: file-storage configure gcp", err)
		return nil, fmt.Errorf("Google Cloud Storage configuration. Please run: file-storage configure gcp")
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
	if keys["GCS_KEY_FILE_PATH"] == "" {
		return nil, fmt.Errorf("GCS service account file key not found. Please run: file-storage configure gcp")
	}
	serviceAccountKeyFile := keys["GCS_KEY_FILE_PATH"]
	client, err := storage.NewClient(*ctx, option.WithCredentialsFile(serviceAccountKeyFile))
	if err != nil {
		return nil, fmt.Errorf("Error: Cannot get GCP client")
	}
	return client, nil
}

func uploadToGCSBucket(bucketName string, filePath string, objectName string) bool {
	ctx := context.Background()
	client, err := getGCSClient(&ctx)
	if err != nil {
		fmt.Println("Failed to authenticate: ", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		fmt.Printf("Failed to upload file: %v", err)
		return false
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Failed to close writer: %v", err)
		return false
	}

	fmt.Printf("File '%s' uploaded to bucket '%s' successfully!\n", objectName, bucketName)
	return true
}
