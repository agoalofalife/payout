package http

import (
	"encoding/json"
	"fmt"
	"github.com/agoalofalife/payout/drivers/yandex"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

var (
	contentTypeDefault = "application/json"
	portDefault        = ":9000"
)

func Start() {
	if port := os.Getenv("PORT"); port == "" {
		port = portDefault

		http.HandleFunc("/", IndexRouterHandler)
		http.HandleFunc("/yandex/balance", YandexBalanceHandler)

		log.Println("Server run, port: " + port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func IndexRouterHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Server is run!"))
}

func YandexBalanceHandler(response http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		panic(err)
	}
	if clientOrderId := request.PostFormValue("clientOrderId"); clientOrderId == "" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Parameter clientOrderId is required."))
	}

	balance := yandex.NewBalance(12)
	balance.Run()
	response.Header().Set("Content-Type", contentTypeDefault)
	fmt.Fprint(response, JsonResponse{"balance": balance.Balance()})
}

type JsonResponse map[string]interface{}

func (jr JsonResponse) String() (str string) {
	byte, err := json.Marshal(jr)
	if err != nil {
		str = ""
		return
	}
	str = string(byte)
	return
}
