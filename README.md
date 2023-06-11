# File Storage CLI Tool

The File Storage CLI Tool is a command-line interface that allows users to configure AWS S3 credentials and perform various operations such as uploading, listing, and deleting files in an S3 bucket.

## Prerequisites

- Go programming language (version 1.16 or later)
- AWS CLI installed and configured with appropriate access credentials
- S3 bucket created in your AWS account

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/file-storage.git
```
2. Navigate to the project directory:
```bash
cd file-storage
```
3. Build the project:
```bash
make build
```
For Windows users, instead of `make build`, run `build.bat`.

## Usage
### Configure AWS Credentials
To configure AWS credentials, run the following command:
```bash
file-storage configure aws
```
This command will prompt you to enter your AWS access key ID, secret access key, and default region. These credentials will be used for subsequent AWS S3 operations.

### Upload a File to S3
To upload a file to an S3 bucket, use the following command:
```bash
file-storage upload aws <REGION> <BUCKET_NAME> <FILE_PATH> <FILE_KEY>
```
Replace `<REGION>` with the desired AWS region, `<BUCKET_NAME>` with the name of your S3 bucket, `<FILE_PATH>` with the local path to the file you want to upload, and `<FILE_KEY>` with the desired object key for the uploaded file.

### List Files in an S3 Bucket  
To list all files in an S3 bucket, execute the following command:
```bash
file-storage ls aws <REGION> <BUCKET_NAME>
```
Replace <REGION> with the AWS region, <BUCKET_NAME> with the name of your S3 bucket, and <FILE_KEY> with the object key of the file you want to delete.  
### Delete a File from an S3 Bucket  
To delete a file from an S3 bucket, use the following command:
```bash
file-storage delete aws <REGION> <BUCKET_NAME> <FILE_KEY>
```
Replace `<REGION>` with the AWS region, `<BUCKET_NAME>` with the name of your S3 bucket, and `<FILE_KEY>` with the object key of the file you want to delete.