package objectstorage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func TestPresignerGeneratePresignedURL(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	mock := newMockMinioClient()
	osClient, err := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	if err != nil {
		t.Errorf("New() failed: %s\n", err)
	}

	svc := osClient.Presigner()
	bucket := "test-bucket"
	object := "test-key"
	expiry := time.Minute
	path := fmt.Sprintf("%s/%s", bucket, object)
	expires := strconv.Itoa(int(expiry.Seconds()))

	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "presigned URL for HTTP GET",
			method: http.MethodGet,
		},
		{
			name:   "presigned URL for HTTP HEAD",
			method: http.MethodHead,
		},
		{
			name:   "presigned URL for HTTP PUT",
			method: http.MethodPut,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.GeneratePresignedURL(context.Background(), tt.method, bucket, object, expiry, url.Values{})
			if err != nil {
				t.Errorf("GeneratePresignedURL() failed: %s\n", err)
			}
			if res.Path != path {
				t.Errorf("URL path mismatch: got %s want %s\n", res.Path, path)
			}
			if exp := res.Query().Get("X-Amz-Expires"); exp != expires {
				t.Errorf("X-Amz-Expires mismatch: got %s want %s\n", exp, expires)
			}
		})
	}
}

func TestPresignerGeneratePresignedURLWithInvalidHTTPMethod(t *testing.T) {
	t.Parallel()

	core := client.NewMgcClient()
	mock := newMockMinioClient()
	osClient, err := New(core, "minioadmin", "minioadmin", WithMinioClientInterface(mock))
	if err != nil {
		t.Errorf("New() failed: %s\n", err)
	}

	svc := osClient.Presigner()
	method := http.MethodDelete
	expected := &InvalidHTTPMethodError{Method: method}

	_, err = svc.GeneratePresignedURL(context.Background(), method, "test-bucket", "test-key", time.Minute, url.Values{})
	if err == nil {
		t.Errorf("GeneratePresignedURL() expected error for HTTP %s method but got nil", method)
	}
	if _, ok := err.(*InvalidHTTPMethodError); !ok {
		t.Errorf("Error mismatch: got %#v want %#v\n", err, expected)
	}
}
