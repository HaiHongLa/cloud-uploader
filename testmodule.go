package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"io"
	"net/http"
	"net/url"
	"os"
)

func testFunc() {
	accountName := "haila"
	accountKey := "7znRPp+dfLmlwnmc7twLrFA02Bn1HAJLaMn54ysERWFG6qv3OHYQQ93nfhzHBJb+01I5v/0TSpnz+AStgYTohA=="
	containerName := "myfiles"

	// Create a shared key credential
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Printf("Failed to create shared key credential: %v\n", err)
		return
	}

	// Create a pipeline and service URL
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	serviceURL := azblob.NewServiceURL(*u, pipeline)
	containerURL := serviceURL.NewContainerURL(containerName)

	// Open the file
	filePath := "D:/UW_Madison/Spring 2023/CS540/homeworks/hw4/hw4-3.zip"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// Create a blob URL and initiate the upload
	blobName := "Homework/hw.zip"
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Detect the content type using http.DetectContentType
	buffer := make([]byte, 512) // Read first 512 bytes to detect content type
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to read file: %v\n", err)
		return
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
		return
	}

	fmt.Printf("Upload completed with status code: %d\n", uploadResponse.Response().StatusCode)
}

