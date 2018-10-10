package http

import (
	"encoding/json"
	"fmt"
	"github.com/agoalofalife/payout/drivers/yandex"
	"log"
	"net/http"
	"os"
)

var (
	port string
	contentTypeDefault  = "application/json"
	portDefault         = ":9000"
	jsonResponseDefault = map[string]interface{}{"result": "", "error": ""}
)

func Start() {
	if port = os.Getenv("PORT"); port == "" {
		port = portDefault
	}
	http.HandleFunc("/", indexRouterHandler)
	http.HandleFunc("/yandex/balance", yandexBalanceHandler)
	http.HandleFunc("/yandex/testDeposition/phone", yandexTestDepositionPhone)

	log.Println("Server run, port: " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexRouterHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Server is run!"))
}

func yandexTestDepositionPhone(response http.ResponseWriter, request *http.Request) {
	var err error
	decoder := json.NewDecoder(request.Body)

	s := struct {
		ClientOrderId int
		Amount float64
		DstAccount int64
	}{}

	err = decoder.Decode(&s)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Parameter Phone is required and expected json."))
	} else {
		testDeposition := yandex.NewTestDeposition(s.ClientOrderId, s.DstAccount, s.Amount, "")
		testDeposition.Run()
		response.Header().Set("Content-Type", contentTypeDefault)
		if testDeposition.IsError() {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}, testDeposition.GetMessageError()))
		} else {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}))
		}
	}
}

// handler api get yandex/balance
func yandexBalanceHandler(response http.ResponseWriter, request *http.Request) {
	var err error
	decoder := json.NewDecoder(request.Body)

	s := struct {
		ClientOrderId int
	}{}
	err = decoder.Decode(&s)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Parameter clientOrderId is required and expected json."))
	} else {
			balance := yandex.NewBalance(s.ClientOrderId)
			balance.Run()
			response.Header().Set("Content-Type", contentTypeDefault)
			if balance.IsError() {
				fmt.Fprint(response, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}, balance.GetMessageError()))
			} else {
				fmt.Fprint(response, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}))
			}
	}
}

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
