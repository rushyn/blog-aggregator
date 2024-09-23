package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rushyn/blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries

}

var apiCfg = apiConfig{}


type payload interface{
	retrunSelf() interface{}
}

type status struct{
	Status string `json:"status"`
}

func (s status) retrunSelf() interface{} {
	return s
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	db, err := sql.Open("postgres", os.Getenv("dbURL"))
	if err != nil {
		log.Fatal(err)
	}

	apiCfg.DB = database.New(db)

	port := os.Getenv("PORT")

	mux := http.NewServeMux()


	mux.HandleFunc("GET /v1/healthz", server_status)
	mux.HandleFunc("GET /v1/err", with_err)
	mux.HandleFunc("POST /v1/users", create_user)
	mux.HandleFunc("GET /v1/users", get_user)
	mux.HandleFunc("GET /users", apiCfg.middAuthenticate(apiCfg.RetrunUserInfo))
	mux.HandleFunc("POST /v1/feeds", apiCfg.middAuthenticate(apiCfg.UpdateFeed))

	

	
	svr := &http.Server{
		Addr: ":" + port,
		Handler: mux,

	}

	
	log.Printf("Http server starting on port: %s\n", port)
	log.Fatal(svr.ListenAndServe())
}

func respondWithJSON(w http.ResponseWriter, code int, p payload){


	data, err := json.Marshal(p.retrunSelf())
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


