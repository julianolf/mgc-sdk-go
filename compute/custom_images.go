package compute

import (
	"context"
	"net/http"

	mgc_http "github.com/MagaluCloud/mgc-sdk-go/internal/http"
)

// Platform represents the system platform.
type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
)

// Architecture represents the system architecure.
type Architecture string

const ArchitectureX86_64 Architecture = "x86/64"

// License indicates if the image software requires a license.
type License string

const (
	LicenseLicensed   License = "licensed"
	LicenseUnlicensed License = "unlicensed"
)

// CreateCustomImageRequest represents the request to create a new custom image.
type CreateCustomImageRequest struct {
	Name         string               `json:"name"`
	Platform     Platform             `json:"platform"`
	Architecture Architecture         `json:"architecture"`
	License      License              `json:"license"`
	URL          string               `json:"url"`
	Requirements *MinimumRequirements `json:"requirements,omitempty"`
	Version      *string              `json:"version,omitempty"`
	Description  *string              `json:"description,omitempty"`
	UEFI         *bool                `json:"uefi,omitempty"`
}

// CustomImageService provides operations for managing custom virtual machine images.
// This interface allows create custom images.
type CustomImageService interface {
	Create(ctx context.Context, req CreateCustomImageRequest) (string, error)
}

// customImageService implements the CustomImageService interface.
// This is an internal implementation that should not be used directly.
type customImageService struct {
	client *VirtualMachineClient
}

// Create creates a new custom image.
// This method makes an HTTP request to publish a new custom image
// and returns the ID of the created image.
func (s *customImageService) Create(ctx context.Context, createReq CreateCustomImageRequest) (string, error) {
	res, err := mgc_http.ExecuteSimpleRequestWithRespBody[struct{ ID string }](
		ctx,
		s.client.newRequest,
		s.client.GetConfig(),
		http.MethodPost,
		"/v1/images/custom",
		createReq,
		nil,
	)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}
