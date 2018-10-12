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
// route /
func indexRouterHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Server is run!"))
}
// route /yandex/balance
func yandexBalanceHandler(response http.ResponseWriter, request *http.Request) {
	var err error
	decoder := json.NewDecoder(request.Body)

	jsonRequest := BaseJsonRequest{}
	err = decoder.Decode(&jsonRequest)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Parameter clientOrderId is required and expected json."))
	} else {
		balance := yandex.NewBalance(jsonRequest.ClientOrderId)
		balance.Run()
		response.Header().Set("Content-Type", contentTypeDefault)
		if balance.IsError() {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}, balance.GetMessageError()))
		} else {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}))
		}
	}
}

// route /yandex/testDeposition/phone
func yandexTestDepositionPhone(response http.ResponseWriter, request *http.Request) {
	var err error
	decoder := json.NewDecoder(request.Body)

	requestJson := newDepositionJsonRequestPhone()

	err = decoder.Decode(&requestJson)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Error json."))
	} else {
		testDeposition := yandex.NewDeposition(yandex.TestDeps, requestJson.ClientOrderId, requestJson.DstAccount, requestJson.Amount, requestJson.Contract)
		testDeposition.Run()
		response.Header().Set("Content-Type", contentTypeDefault)
		if testDeposition.IsError() {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}, testDeposition.GetMessageError()))
		} else {
			fmt.Fprint(response, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}))
		}
	}
}
