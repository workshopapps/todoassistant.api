package awsService

import (
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	awss3 "test-va/internals/amazon/awsS3"
)

func TestUploadImage(t *testing.T) {
	// create an awsSrv struct with the mock s3session
	s3session, err := awss3.NewAWSSession("AKIA5G6AVJRXO42IGWER", "9C4gv28jjxk2dBxkl49+FO/SUksaGEnpYRncn5iq", "")
	if err != nil {
		log.Println("Error Connecting to AWS S3: ", err)
		t.Errorf("Unexpected error: %v", err)
	}

	a := NewAWSSrv(s3session)
	// create a mock file to be uploaded
	mockFile := CreateFormFile()

	// call the UploadImage function
	err = a.UploadImage(mockFile, "test_image.jpg")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func CreateFormFile() *multipart.FileHeader {
	// create a mock file
	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "mock_image.jpg")
	mockFile, err := os.Open(path)
	// fmt.Println(path)
	if err != nil {
		panic(err)
	}
	defer mockFile.Close()

	// create header
	header := make(map[string][]string)
	header["Content-Disposition"] = []string{`form-data; name="file"; filename="mock_image.jpg"`}
	header["Content-Type"] = []string{"image/jpeg"}

	// create a form file header
	fileHeader := &multipart.FileHeader{
		Filename: path,
		Header:   header,
		Size:     1289231,
	}

	// create a form
	// form := &multipart.Form{
	// 	File: map[string][]*multipart.FileHeader{
	// 		"file": {fileHeader},
	// 	},
	// }

	return fileHeader
}
