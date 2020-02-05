package main

import (
    "github.com/minio/minio-go/v6"
    "log"
)

func main() {
    endpoint := "play.min.io"
    accessKeyID := "Q3AM3UQ867SPQQA43P2F"
    secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
    useSSL := true

    // Initialize minio client object.
    minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
    if err != nil {
        log.Fatalln(err)
    }

    // Make a new bucket called mymusic.
    bucketName := "mymusic"
    location := "us-east-1"

    err = minioClient.MakeBucket(bucketName, location)
    if err != nil {
        // Check to see if we already own this bucket (which happens if you run this twice)
        exists, errBucketExists := minioClient.BucketExists(bucketName)
        if errBucketExists == nil && exists {
            log.Printf("We already own %s\n", bucketName)
        } else {
            log.Fatalln(err)
        }
    } else {
        log.Printf("Successfully created %s\n", bucketName)
    }

    // Upload the zip file
    objectName := "golden-oldies.zip"
    filePath := "/tmp/golden-oldies.zip"
    contentType := "application/zip"

    // Upload the zip file with FPutObject
    n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType:contentType})
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}