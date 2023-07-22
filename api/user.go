package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id string `json:"id"`
}

type ConnectJob struct {
	JobType string `json:"type"`
	Id      string `json:"id"`
	Links   JobLink
}

type JobLink struct {
	Self string `json:"self"`
}

func CreateUser(accessToken string) string {
	createUserUrl := BaseUrl + "users"
	createUserBody := []byte(`{
		"email": "gavin@hooli.com",
		"mobile": "+61410888666",
		"firstName": "Joe",
		"lastName": "Bloggs"
	}`)

	r, err := http.NewRequest("POST", createUserUrl, bytes.NewBuffer(createUserBody))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	payload := &User{}
	derr := json.NewDecoder(res.Body).Decode(payload)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusCreated {
		panic(res.Status)
	}

	fmt.Println("User created!")
	return payload.Id
}

func ConnectUser(accessToken string, userId string) string {
	const baseURL = "https://au-api.basiq.io/"

	connUserUrl := baseURL + "users/" + userId + "/connections"
	connUserBody := []byte(`{
		"loginId": "gavinBelson",
		"password": "hooli2016",
		"institution": {
			"id": "AU00000"
		}
	}`)

	r, err := http.NewRequest("POST", connUserUrl, bytes.NewBuffer(connUserBody))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	payload := &ConnectJob{}
	derr := json.NewDecoder(res.Body).Decode(payload)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusAccepted {
		panic(res.Status)
	}

	return payload.Id
}
