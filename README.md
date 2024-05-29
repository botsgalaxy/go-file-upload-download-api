### Go File Upload/Download Service

### This is a simple Go application that provides a file upload and download service. It allows users to upload files to the server and retrieve them via direct URLs.

### Features

    - Upload files with a maximum size limit of 300 MB.
    - Generate direct download URLs for uploaded files.
    - Automatic removal of uploaded files after 30 minutes.

### Prerequisites

    Go 1.16 or higher installed on your system.
    Gin web framework for Go. You can install it using:
    go get -u github.com/gin-gonic/gin

## Usage

- Clone this repository to your local machine:
- git clone https://github.com/botsgalaxy/go-file-upload-download-api.git
- Navigate to the project directory:
- cd go-file-upload-download
- Run the application:
- go run main.go
- The server will start running at http://localhost:8080.

## Endpoints

### Upload Endpoint

    URL: /upload
    Method: POST
    Request Parameters:
        file: The file to be uploaded.
    Response:
        On success:

        json

            ```{
            "status": "success",
            "data": {
            "url": "/download/file_id"
            }
            }```

        On failure:

            json

        ```{
            "status": "error",
            "message": "Error message here"
        }```

### Download Endpoint

    URL: /download/:file_id
    Method: GET
    Response: The file for download.

### Configuration

You can modify the following configurations in the main.go file:

- MaxFileSize: Maximum allowed file size in bytes.
- MaxUploadTimeout: Timeout for file removal after upload (in minutes).
- UploadsDir: Directory path where uploaded files will be stored.

License

This project is licensed under the MIT License - see the LICENSE file for details.
