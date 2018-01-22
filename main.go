package main

import (
	"fmt"
	"log"
	"os"

	minio "github.com/minio/minio-go"
)

func recreateBucket(minioClient *minio.Client) {
	list(minioClient, func(object minio.ObjectInfo) {
		fmt.Println("deleting", object.Key)

		err := minioClient.RemoveObject("mybucket", object.Key)
		if err != nil {
			fmt.Println(err)
		}
	})

	err := minioClient.RemoveBucket("mybucket")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = minioClient.MakeBucket("mybucket", "us-east-1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully created bucket")
}

func put(minioClient *minio.Client) {
	file, err := os.Open("Makefile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := minioClient.PutObject("mybucket", "myobject", file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		UserMetadata: map[string]string{
			"hhhhhh": "helloworld",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded bytes: ", n)
}

func list(minioClient *minio.Client, do func(minio.ObjectInfo)) {
	// Create a done channel to control 'ListObjectsV2' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	isRecursive := true
	objectCh := minioClient.ListObjectsV2("mybucket", "", isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		do(object)
	}
}

func statobject(minioClient *minio.Client) {
	i, err := minioClient.StatObject("mybucket", "myobject", minio.StatObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("stat %+v\n", i)
}

func main() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "secretaccesskey"
	secretAccessKey := "password!!"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	recreateBucket(minioClient)
	put(minioClient)
	list(minioClient, func(object minio.ObjectInfo) {
		fmt.Printf("%+v\n", object)
	})
	statobject(minioClient)
}
