package kideungeo

import (
	"net/url"
	"path"
)

func JoinURL(baseURL string, paths ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, err
	}
	newPaths := append([]string{u.Path}, paths...)
	u.Path = path.Join(newPaths...)
	return u.String(), nil
}
