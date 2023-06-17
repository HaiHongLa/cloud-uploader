# Cloud Storage CLI Tool

The File Storage CLI Tool is a command-line interface that allows users to configure credentials and perform various operations such as downloading, uploading, listing, and deleting files in AWS S3 buckets or Microsoft Azure Blob Storage containers.

- [AWS S3](#for-aws-s3)
- [Microsoft Azure Blob Storage](#for-microsoft-azure-blob-storage)
## Prerequisites

- Go 1.16 or later

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
## For AWS S3
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
### Get a File from S3 Bucket to Local Machine
To download a file from an S3 bucket to your local machine, use the following command:
```bash
file-storage get aws <REGION> <BUCKET_NAME> <FILE_KEY> <FILE_PATH>
```
Replace `<REGION>` with the desired AWS region, `<BUCKET_NAME>` with the name of your S3 bucket, `<FILE_KEY>` with the object key of the file you want to download, and `<FILE_PATH>` with the local path where you want to save the downloaded file.
## For Microsoft Azure Blob Storage
### Configure Azure Credentials
To configure Azure credentials, run the following command:  
```bash
file-storage configure azure
```
This command will prompt you to enter your Azure Storage account name and access key. These credentials will be used for subsequent Azure Blob Storage operations.  
### Upload a File to Azure Blob Storage
To upload a file to an Azure Blob Storage container, use the following command:
```bash
file-storage upload azure <CONTAINER_NAME> <FILE_PATH> <BLOB_NAME>
```
Replace `<CONTAINER_NAME>` with the name of your Azure Blob Storage container, `<FILE_PATH>` with the local path to the file you want to upload, and `<BLOB_NAME>` with the desired blob name for the uploaded file.

### List Files in an Azure Blob Storage Container
To list all files in an Azure Blob Storage container, execute the following command:
```bash
file-storage ls azure <CONTAINER_NAME>
```
Replace `<CONTAINER_NAME>` with the name of your Azure Blob Storage container.

### Delete a Blob from an Azure Blob Storage Container
To delete a blob from an Azure Blob Storage container, use the following command:
```bash
file-storage delete azure <CONTAINER_NAME> <BLOB_NAME>
```
Replace `<CONTAINER_NAME>` with the name of your Azure Blob Storage container and `<BLOB_NAME>` with the name of the blob you want to delete.

### Get a Blob from an Azure Blob Storage Container to Local Machine
To download a blob from an Azure Blob Storage container to your local machine, use the following command:
```bash
file-storage get azure <CONTAINER_NAME> <BLOB_NAME> <FILE_PATH>
```
Replace `<CONTAINER_NAME>` with the name of your Azure Blob Storage container, `<BLOB_NAME>` with the name of the blob you want to download, and `<FILE_PATH>` with the local path where you want to save the downloaded file.

Remember to replace `<CONTAINER_NAME>`, `<BLOB_NAME>`, and `<FILE_PATH>` with the actual names and paths specific to your scenario.

With these commands, you can manage your Azure Blob Storage containers and perform operations such as uploading, listing, deleting, and downloading files.