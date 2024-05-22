package utils

import (
	"fmt"
	"net/url"
)

func GenerateUrl(domain string, shortID string) string {

	return fmt.Sprintf("http://%s/%s", domain, shortID)

}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// If the URL has a non-empty scheme and host, it is a valid URL.
	return u.Scheme != "" && u.Host != ""
}
