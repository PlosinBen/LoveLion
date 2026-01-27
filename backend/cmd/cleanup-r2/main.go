package main

import (
	"context"
	"fmt"
	"log"

	"lovelion/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	log.Println("🧹 Starting R2 Bucket Cleanup...")

	cfg := config.Load()
	if cfg.R2Bucket == "" {
		log.Println("⚠️  Skipping R2 cleanup: R2_BUCKET_NAME not set")
		return
	}
	if cfg.R2AccountID == "" || cfg.R2AccessKey == "" || cfg.R2SecretKey == "" {
		log.Println("⚠️  Skipping R2 cleanup: Missing R2 credentials")
		return
	}

	log.Printf("Target Bucket: %s", cfg.R2Bucket)

	// Setup S3 Client for R2
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKey, cfg.R2SecretKey, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	client := s3.NewFromConfig(awsCfg)

	// List and delete loop
	var deletedCount int
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.R2Bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("❌ Failed to list objects: %v", err)
		}

		if len(page.Contents) == 0 {
			continue
		}

		var objects []types.ObjectIdentifier
		for _, obj := range page.Contents {
			objects = append(objects, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		// Delete objects in batches (max 1000 per request, handled by loop implicitly if page size is default)
		_, err = client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(cfg.R2Bucket),
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			log.Fatalf("❌ Failed to delete objects: %v", err)
		}

		deletedCount += len(objects)
		log.Printf("Deleted %d objects...", deletedCount)
	}

	log.Printf("✅ R2 Cleanup complete! Total deleted: %d", deletedCount)
}
