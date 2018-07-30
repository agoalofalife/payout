package http

import (
	"github.com/agoalofalife/payout/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	templateExpectedResponse = "handler returned unexpected body: got %v want %v"
	templateWongStatusCode   = "handler returned wrong status code: got %v want %v"
)

func TestIndex(t *testing.T) {
	request := utils.FakeRequest("/", "GET", t)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexRouterHandler)
	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf(templateWongStatusCode,
			status, http.StatusOK)
	}

	expected := `Server is run!`
	if responseRecorder.Body.String() != expected {
		t.Errorf(templateExpectedResponse,
			responseRecorder.Body.String(), expected)
	}
}
