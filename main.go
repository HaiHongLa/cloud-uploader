package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args
	fmt.Println(args)
	if len(args) <= 1 {
		fmt.Println("Usage: configure")
		return
	} else if args[1] == "configure" {
		if !configHandler(args[2]) {
			fmt.Println("An error occured when configuring AWS credentials")
			return
		}
	} else if args[1] == "upload" {
		if !uploadHandler(args) {
			fmt.Println("An error occured when uploading")
		}
	}

	fmt.Println("Usage: configure")
	return
}

func configHandler(platform string) bool {
	if strings.ToLower(platform) == "aws" {
		return aws_config()
	}
	return false
}

func uploadHandler(args []string) bool {
	if strings.ToLower(args[2]) == "aws" {
		if len(args) != 7 {
			fmt.Println("Usage: file-storage upload aws <REGION> <BUCKET_NAME> <FILE_PATH> <FILE_KEY>")
			return false
		}
		return aws_upload(args[3], args[4], args[5], args[6])
	}
	return false
}