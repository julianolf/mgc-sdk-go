package objectstorage

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// PresignedService provides operations for generating presigned URLs.
type PresignedService interface {
	GeneratePresignedURL(ctx context.Context, method string, bucketName string, objectKey string, expiry time.Duration, reqParams url.Values) (*url.URL, error)
}

// presignedService implements the PresignedService interface.
type presignedService struct {
	client *ObjectStorageClient
}

// Generates a presigned URL for the given HTTP method operations.
func (p *presignedService) GeneratePresignedURL(ctx context.Context, method string, bucketName string, objectKey string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	switch method {
	case http.MethodGet:
		return p.client.minioClient.PresignedGetObject(ctx, bucketName, objectKey, expiry, reqParams)
	case http.MethodHead:
		return p.client.minioClient.PresignedHeadObject(ctx, bucketName, objectKey, expiry, reqParams)
	case http.MethodPut:
		return p.client.minioClient.PresignedPutObject(ctx, bucketName, objectKey, expiry)
	default:
		return nil, &InvalidHTTPMethodError{Method: method}
	}
}
