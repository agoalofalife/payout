package utils

import (
	"net/http"
	"testing"
)

func FakeRequest(path string, method string, t *testing.T) *http.Request {
	request, err := http.NewRequest(method, path, nil)
	// (*Request, error)
	if err != nil {
		t.Fatal(err)
	}
	return request
}
