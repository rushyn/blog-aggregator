package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)


func get_user (w http.ResponseWriter, r *http.Request) {


	ApiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
	
	ctx := context.Background()
	databaseUser, err := apiCfg.DB.GetUser(ctx, ApiKey)
	if err != nil{
		log.Println(err)
		return
	}

	user := user{
		ID: databaseUser.ID,
		Created_at: databaseUser.CreatedAt,
		Updated_at: databaseUser.UpdatedAt,
		Name: databaseUser.Name,
		Api_key: databaseUser.ApiKey,
	}


	respondWithJSON(w, 200, user)

}