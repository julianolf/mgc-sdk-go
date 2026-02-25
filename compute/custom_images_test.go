package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomImageService_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		req        CreateCustomImageRequest
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name: "successful creation",
			req: CreateCustomImageRequest{
				Name:         "test-image",
				Platform:     PlatformLinux,
				Architecture: ArchitectureX86_64,
				License:      LicenseUnlicensed,
				URL:          "https://br-se1.magaluobjects.com/bucket/image.qcow2",
			},
			response:   `{"id": "8cf5c6d9-d5c5-4af9-bd1b-c17d032dc761"}`,
			statusCode: http.StatusOK,
			wantID:     "8cf5c6d9-d5c5-4af9-bd1b-c17d032dc761",
			wantErr:    false,
		},
		{
			name: "empty name",
			req: CreateCustomImageRequest{
				Platform:     PlatformLinux,
				Architecture: ArchitectureX86_64,
				License:      LicenseUnlicensed,
				URL:          "https://br-se1.magaluobjects.com/bucket/image.qcow2",
			},
			response:   `{"error": "name is required"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "invalid architecture",
			req: CreateCustomImageRequest{
				Name:         "test-image",
				Platform:     PlatformLinux,
				Architecture: Architecture("arm64"),
				License:      LicenseUnlicensed,
				URL:          "https://br-se1.magaluobjects.com/bucket/image.qcow2",
			},
			response:   `{"error": "invalid architecture"}`,
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "server error",
			req: CreateCustomImageRequest{
				Name:         "test-image",
				Platform:     PlatformLinux,
				Architecture: ArchitectureX86_64,
				License:      LicenseUnlicensed,
				URL:          "https://br-se1.magaluobjects.com/bucket/image.qcow2",
			},
			response:   `{"error": "internal error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "duplicate name",
			req: CreateCustomImageRequest{
				Name:         "test-duplicated-image",
				Platform:     PlatformLinux,
				Architecture: ArchitectureX86_64,
				License:      LicenseUnlicensed,
				URL:          "https://br-se1.magaluobjects.com/bucket/image.qcow2",
			},
			response:   `{"error": "image name already exists"}`,
			statusCode: http.StatusConflict,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := testClient(server.URL)
			gotID, err := client.CustomImages().Create(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("Create() got = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
