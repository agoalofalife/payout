package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	templateExpectedResponse = "handler returned unexpected body: got %v want %v"
	templateWongStatusCode   = "handler returned wrong status code: got %v want %v"
)

func TestIndex(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(StartRouterHandler)
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
