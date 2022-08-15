package util

import (
	"net/url"
)

func JoinPath(base string, elem string) (string, error) {
	bu, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(elem)
	if err != nil {
		return "", nil
	}
	return bu.ResolveReference(u).String(), nil
}
