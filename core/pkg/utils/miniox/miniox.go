package miniox

import (
	"fmt"
	"strings"
)

func FulfillURL(externalURL string, bucketName string, relativePath string) string {
	if relativePath == "" {
		return ""
	}

	if strings.HasPrefix(relativePath, "http") {
		return relativePath
	}

	url := strings.TrimSuffix(externalURL, "/")
	return fmt.Sprintf("%s/%s/%s", url, bucketName, relativePath)
}
