package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	s3Bucket string
)

func init() {
	// Config aws sdk
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				"",
				"",
				"")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}

	// Create an Amazon S3 service client
	s3Client = s3.NewFromConfig(cfg)
	s3Bucket = "go-expert"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	wg := sync.WaitGroup{}
	uploadControl := make(chan struct{}, 10)
	errorFileUpdate := make(chan string, 5)

	go func() {
		for {
			select {
			case file := <-errorFileUpdate:
				uploadControl <- struct{}{}
				wg.Add(1)
				go uploadFile(file, uploadControl, errorFileUpdate, &wg)
			}
		}
	}()

	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}
		uploadControl <- struct{}{}
		wg.Add(1)
		go uploadFile("./tmp/" + files[0].Name(), uploadControl, errorFileUpdate, &wg)
	}
	wg.Wait()

}

func uploadFile(fileName string, uploading chan<- struct{}, retry chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	uploading <- struct{}{}
	err := uploadFileToS3(fileName)
	if err != nil {
		fmt.Printf("Error uploading file %q: %v, retrying......\n", fileName, err)
		retry <- fileName
		return
	}
	err = os.Remove(fileName)
	if err != nil {
		fmt.Printf("Error removing file %q: %v\n", fileName, err)
	}
}

func uploadFileToS3(fileName string) error {
	file, err := os.Open(fileName)
	fmt.Printf("Uploading file %q to S3 bucket %q...\n", fileName, s3Bucket)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}
