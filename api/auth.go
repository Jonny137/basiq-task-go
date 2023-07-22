package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type AuthPost struct {
	AccessToken string `json:"access_token"`
}

func AuthorizeUser() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	authUrl := BaseUrl + "token"

	r, err := http.NewRequest("POST", authUrl, nil)
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic "+apiKey)
	r.Header.Add("basiq-version", "2.1")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	payload := &AuthPost{}
	derr := json.NewDecoder(res.Body).Decode(payload)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	fmt.Println("Authorization successful!")
	return payload.AccessToken
}
