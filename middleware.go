package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rushyn/blog-aggregator/internal/database"
)


type returnFeedUpdate struct{
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (r returnFeedUpdate) retrunSelf() interface{} {
	return r
}


func dbFeedUpdateToReturnFeedUpdate(c database.Feed) returnFeedUpdate {
	return returnFeedUpdate{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
		Url:       c.Url,
		UserID:    c.UserID,
	}
}



type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middAuthenticate(handler authedHandler)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		ApiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")

		if ApiKey == ""{
			respondWithError(w, 401, "no valid key")
			return
		}

		ctx := context.Background()
		databaseUser, err := apiCfg.DB.GetUser(ctx, ApiKey)
		if err != nil{
			log.Println(err)
			respondWithError(w, 403, err.Error())
			return
		}

		handler(w, r, databaseUser)
	}
}

func (cfg *apiConfig) RetrunUserInfo(w http.ResponseWriter, r *http.Request, user database.User){

	respondWithJSON(w, 200, dbUserToUser(user))
}

func (cfg *apiConfig) UpdateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	
	type FeedUpdate struct{
		Name string `json:"name"`
		Url string `json:"url"`
	}

	update := FeedUpdate{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&update)
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

	

	ctx := context.Background()
	newFeed, err := apiCfg.DB.CreateFeedEntry(ctx,
	database.CreateFeedEntryParams{
		ID:        newUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      update.Name,
		Url:       update.Url,
		UserID:    user.ID,
	})
	if err != nil{
		log.Println(err)
		respondWithError(w, 500, "fail to update database")
		return
	}

	respondWithJSON(w, 200, dbFeedUpdateToReturnFeedUpdate(newFeed))

}