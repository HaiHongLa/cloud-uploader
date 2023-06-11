package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("For instructions run: file-storage help")
		return
	} else if args[1] == "configure" {
		if !configHandler(args[2]) {
			fmt.Println("An error occured when configuring AWS credentials")
			return
		}
	} else if args[1] == "upload" {
		if !uploadHandler(args) {
			fmt.Println("An error occured when uploading")
			return
		}
		return
	} else if args[1] == "ls" {
		if !lsHandler(args) {
			fmt.Println("An error occured when trying to list files")
			return
		}
		return
	} else if args[1] == "delete" {
		if !deleteHandler(args) {
			fmt.Println("An error occured when trying to delete file")
			return
		}
		return
	}
	fmt.Println("Command not supported")
	return
}

func configHandler(platform string) bool {
	if strings.ToLower(platform) == "aws" {
		return configureAWSS3()
	}
	return false
}

func uploadHandler(args []string) bool {
	fmt.Println(args)
	if strings.ToLower(args[2]) == "aws" {
		if len(args) != 7 {
			fmt.Println("Usage: file-storage upload aws <REGION> <BUCKET_NAME> <FILE_PATH> <FILE_KEY>")
			return false
		}
		return uploadToAWSS3Bucket(args[3], args[4], args[5], args[6])
	}
	return false
}

func lsHandler(args []string) bool {
	if strings.ToLower(args[2]) == "aws" {
		if len(args) != 5 {
			fmt.Println("Usage: file-storage ls aws <REGION> <BUCKET_NAME>")
			return false
		}
		return listAllFilesFromAWSS3Bucket(args[3], args[4])
	}
	return false
}

func deleteHandler(args []string) bool {
	if strings.ToLower(args[2]) == "aws" {
		if len(args) != 6 {
			fmt.Println("Usage: file-storage delete aws <REGION> <BUCKET_NAME> <FILE_KEY>")
			return false
		}
		return deleteFromS3Bucket(args[3], args[4], args[5])
	}
	return false
}