package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


type payload interface{
	message() status
}

type status struct{
	Status string `json:"status"`
}

func (self status) message() status {
	return self
}



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", server_status)
	mux.HandleFunc("GET /v1/err", with_err)


	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Http server starting on port: %s\n", port)
	log.Fatal(svr.ListenAndServe())
}


func respondWithJSON(w http.ResponseWriter, code int, p payload){


	data, err := json.Marshal(p.message())
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	
	type error struct {
		Error string `json:"error"`
	}
	
	errorReturn := error{
		Error: msg,
	}

	data, err := json.Marshal(errorReturn)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}


func server_status (w http.ResponseWriter, r *http.Request) {
	

	status := status {
		Status: "ok",
	}
	
	respondWithJSON(w, 200, status)

}

func with_err (w http.ResponseWriter, r *http.Request) {
	
	respondWithError(w, 500, "Internal Server Error")

}