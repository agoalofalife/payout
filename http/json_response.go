package http

import "encoding/json"

type JsonResponse map[string]interface{}
// function helper
// pass result and error return type JsonResponse
func newJsonResponse(result map[string]interface{}, error ...string) JsonResponse {
	return JsonResponse{
		"error":  error,
		"result": result,
	}
}
// conversion struct in json string
func (jr JsonResponse) String() (str string) {
	byte, err := json.Marshal(jr)
	if err != nil {
		str = ""
		return
	}
	str = string(byte)
	return
}

