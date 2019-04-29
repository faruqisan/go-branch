package branch

import (
	"context"
	"fmt"

	"github.com/faruqisan/go-branch/httpclient"
)

type (
	// CreateLinkRequest struct define request parameter to create link
	CreateLinkRequest struct {
		BranchKey string   `json:"branch_key"`
		LinkData  LinkData `json:"data"`
	}

	// CreateLinkResponse struct define response data from create link
	CreateLinkResponse struct {
		URL string `json:"url"`
	}

	// LinkData struct define link data request, the field is omitted when empty
	LinkData struct {
		DesktopURL          string `json:"$desktop_url,omitempty"`
		AndroidURL          string `json:"$android_url,omitempty"` // deeplink
		IOSURL              string `json:"$ios_url,omitmepty"`     // deeplink
		AndroidDeeplinkPath string `json:"$android_deeplink_path,omitempty"`
	}
)

var (
	createLinkURL = fmt.Sprintf("%s/%s/url", baseURL, version)
)

// CreateLink function will crete the branch link
func (e *Engine) CreateLink(ctx context.Context, req CreateLinkRequest) (string, error) {
	var (
		err error
		r   CreateLinkResponse
	)

	_, err = e.client.PostJSON(ctx, createLinkURL, httpclient.JSONHeader, req, &r)
	if err != nil {
		return "", err
	}

	return r.URL, err
}
