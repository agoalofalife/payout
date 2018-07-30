package http

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

func Start() {
	if port := os.Getenv("PORT"); port == "" {
		port = ":9000"
		http.HandleFunc("/", StartRouterHandler)

		log.Println("Server run, port: " + port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func StartRouterHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Server is run!"))
}
