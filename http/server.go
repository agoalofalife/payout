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
	http.HandleFunc("/yandex/makeDeposition/phone", yandexMakeDepositionPhone)
	http.HandleFunc("/yandex/testDeposition/purse", yandexTestDepositionPurse)

	log.Println("Server run, port: " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
// route /
func indexRouterHandler(res http.ResponseWriter, request *http.Request) {
	res.Write([]byte("Server is run!"))
}
// route /yandex/balance
func yandexBalanceHandler(res http.ResponseWriter, req *http.Request) {
	var err error
	decoder := json.NewDecoder(req.Body)

	jsonRequest := BaseJsonRequest{}
	err = decoder.Decode(&jsonRequest)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Parameter clientOrderId is required and expected json."))
	} else {
		balance := yandex.NewBalance(jsonRequest.ClientOrderId)
		balance.Run()
		res.Header().Set("Content-Type", contentTypeDefault)
		if balance.IsError() {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}, balance.GetMessageError()))
		} else {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"balance": balance.Balance()}))
		}
	}
}

// route /yandex/testDeposition/phone
func yandexTestDepositionPhone(res http.ResponseWriter, req *http.Request) {
	wrapDepositionPhone(res, req, yandex.TestDeps)
}
// route /yandex/makeDeposition/phone
func yandexMakeDepositionPhone(res http.ResponseWriter, req *http.Request)  {
	wrapDepositionPhone(res, req, yandex.MakeDeps)
}

// wrapper phone deposition for make and test
func wrapDepositionPhone(res http.ResponseWriter, req *http.Request, deposition yandex.TypeDeposition)  {
	var err error
	decoder := json.NewDecoder(req.Body)

	requestJson := newDepositionJsonRequestPhone()

	err = decoder.Decode(&requestJson)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error json."))
	} else {
		testDeposition := yandex.NewDeposition(deposition, requestJson.ClientOrderId, requestJson.DstAccount, requestJson.Amount, requestJson.Contract)
		testDeposition.Run()
		res.Header().Set("Content-Type", contentTypeDefault)
		if testDeposition.IsError() {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}, testDeposition.GetMessageError()))
		} else {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}))
		}
	}
}

func yandexTestDepositionPurse(res http.ResponseWriter, req *http.Request)  {
	wrapDepositionPurse(res, req, yandex.TestDeps)
}

// wrapper purse deposition for make and test
func wrapDepositionPurse(res http.ResponseWriter, req *http.Request, deposition yandex.TypeDeposition)  {
	var err error
	decoder := json.NewDecoder(req.Body)

	requestJson := newDepositionJsonRequestPhone()

	err = decoder.Decode(&requestJson)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error json."))
	} else {
		testDeposition := yandex.NewDeposition(deposition, requestJson.ClientOrderId, requestJson.DstAccount, requestJson.Amount, requestJson.Contract)
		testDeposition.Run()
		res.Header().Set("Content-Type", contentTypeDefault)
		if testDeposition.IsError() {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}, testDeposition.GetMessageError()))
		} else {
			fmt.Fprint(res, newJsonResponse(map[string]interface{}{"success": testDeposition.IsSuccess()}))
		}
	}
}