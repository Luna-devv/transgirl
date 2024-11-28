package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func fetchFileNames() (int, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to load AWS config: %w", err)
	}

	ctx := context.TODO()
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(awsBucket),
		Prefix: aws.String(prefix),
	})

	index := 0
	mapMutex.Lock()
	defer mapMutex.Unlock()

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch bucket objects: %w", err)
		}

		for _, object := range page.Contents {
			fileMap[index] = *object.Key
			index++
		}
	}

	log.Printf("Fetched %d files from bucket %s\n", len(fileMap), awsBucket)
	return len(fileMap), nil
}
