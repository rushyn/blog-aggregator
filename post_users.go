package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rushyn/blog-aggregator/internal/database"
)



type user struct{
	ID uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Name string `json:"name"`
	Api_key string `json:"api_key"`
}

func (u user) retrunSelf() interface{} {
	return u
}

func (u *user) populateSelf(c database.CreateUserParams) {
	u.ID = c.ID
	u.Created_at = c.CreatedAt
	u.Updated_at = c.UpdatedAt
	u.Name = c.Name
	u.Api_key = c.ApiKey
}



func create_user (w http.ResponseWriter, r *http.Request) {
	
	type userName struct{
		Name string `json:"name"`
	}

	name := userName{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&name)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error getting uuid: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	keyLength := 32
	apiByte := make([]byte, keyLength)
	_, err = rand.Read(apiByte)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	apiKey := hex.EncodeToString(apiByte)

	databaseEntrey := database.CreateUserParams{
		ID: newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name.Name,
		ApiKey: apiKey,
	}

	ctx := context.Background()
	insertedAuthor, err := apiCfg.DB.CreateUser(ctx, databaseEntrey)
	if err != nil{
		log.Println(err)
		return
	}

	log.Println(insertedAuthor)

	newUser := user{}
	newUser.populateSelf(databaseEntrey)
	respondWithJSON(w, 200, newUser)

}