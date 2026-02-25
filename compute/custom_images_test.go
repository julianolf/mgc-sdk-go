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

func TestCustomImageService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful request",
			id:   "86a304b0-dc28-454e-9448-5275c4008dfa",
			response: `{
				 "id": "86a304b0-dc28-454e-9448-5275c4008dfa",
				 "name": "test",
				 "status": "active",
				 "platform": "linux",
				 "license": "unlicensed",
				 "requirements": {
				  "vcpu": 1,
				  "ram": 1,
				  "disk": 3
				 },
				 "version": "1.0.0",
				 "description": "Test",
				 "metadata": {
				  "uefi": "true"
				 }
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "image not found",
			id:         "a0db5832-3767-4335-8a89-9b46ce636790",
			response:   `{"message": "Image not found"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "server error",
			id:         "86a304b0-dc28-454e-9448-5275c4008dfa",
			response:   `{"message": "Internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(tt.statusCode)
						w.Write([]byte(tt.response))
					},
				),
			)
			defer server.Close()

			client := testClient(server.URL)
			got, err := client.CustomImages().Get(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Get() expected erro, got nil")
					return
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error: %v", err)
					return
				}
				if got.ID != tt.id {
					t.Errorf("Get() got ID %s, want %s", got.ID, tt.id)
				}
			}
		})
	}
}
