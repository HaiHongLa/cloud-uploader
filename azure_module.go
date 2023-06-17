package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func configureAzureBlobStorage() bool {
	fmt.Print("Enter account name: ")
	var accountName string
	fmt.Scanln(&accountName)
	if len(accountName) == 0 {
		fmt.Println("Please enter your account name")
		return false
	}
	fmt.Print("Enter account access key: ")
	var accessKey string
	fmt.Scanln(&accessKey)
	if len(accessKey) == 0 {
		fmt.Println("Please enter your account access key")
		return false
	}
	err := os.MkdirAll("cloudConfig", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return false
	}
	azureConfig := "AZURE_ACCOUNT_NAME=" + accountName + "\nAZURE_ACCESS_KEY=" + accessKey
	err = ioutil.WriteFile("cloudconfig/AzureConfig", []byte(azureConfig), 0644)
	if err != nil {
		fmt.Println("Error creating file", err)
		return false
	}
	fmt.Println("Configured Microsoft Azure successfully")
	return true
}

func createServiceURL() (*azblob.ServiceURL, error) {
	configFile, err := os.Open("cloudConfig/AzureConfig")
	if err != nil {
		return nil, fmt.Errorf("Azure Configuration not found. Configure your Azure credentials by running: file-storage configure azure. Error: %v", err)
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	keys := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			continue
		}
		keys[pair[0]] = pair[1]
	}
	if keys["AZURE_ACCOUNT_NAME"] == "" {
		return nil, fmt.Errorf("Azure account name is missing. Please run: file-storage configure azure")
	}
	if keys["AZURE_ACCESS_KEY"] == "" {
		return nil, fmt.Errorf("Azure account access key is missing. Please run: file-storage configure azure")
	}
	accountName := keys["AZURE_ACCOUNT_NAME"]
	accountKey := keys["AZURE_ACCESS_KEY"]
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %v", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	serviceURL := azblob.NewServiceURL(*u, pipeline)
	return &serviceURL, nil
}

func uploadToAzureBlobStorage(containerName string, filePath string, blobName string) bool {
	serviceURL, err := createServiceURL()
	if err != nil {
		fmt.Println("Failed to create service URL", err)
		return false
	}
	containerURL := serviceURL.NewContainerURL(containerName)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return false
	}
	defer file.Close()
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Detect the content type using http.DetectContentType
	buffer := make([]byte, 512) // Read first 512 bytes to detect content type
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to read file: %v\n", err)
		return false
	}

	contentType := http.DetectContentType(buffer)

	uploadResponse, err := azblob.UploadFileToBlockBlob(context.TODO(), file, blobURL,
		azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: contentType,
			},
		})

	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return false
	}

	fmt.Printf("Upload completed with status code: %d\n", uploadResponse.Response().StatusCode)
	return true
}

func listAllFilesFromBlobContainer(containerName string) bool {
	serviceURL, err := createServiceURL()
	if err != nil {
		fmt.Println("Failed to create service URL", err)
		return false
	}
	containerURL := serviceURL.NewContainerURL(containerName)
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(context.TODO(), marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			fmt.Println("Error: Cannot list blobs")
			return false
		}

		// Process each blob in the segment
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Println(blobInfo.Name)
		}

		// Update the marker to retrieve the next segment of blobs
		marker = listBlob.NextMarker
	}
	return true
}

func deleteBlobFromContainer(containerName string, blobName string) bool {
	serviceURL, err := createServiceURL()
	if err != nil {
		fmt.Println("Failed to create service URL", err)
		return false
	}
	containerURL := serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlobURL(blobName)
	_, err = blobURL.Delete(context.TODO(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		fmt.Printf("failed to delete blob %v\n", blobName)
		fmt.Println(err)
		return false
	}
	fmt.Printf("Successfully deleted blob %v\n", blobName)
	return true
}