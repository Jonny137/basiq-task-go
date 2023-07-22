package api

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Step struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Result *Result
}

type Steps struct {
	Steps []Step
}

func CheckConnectionJob(accessToken string, jobId string) string {
	jobUrl := BaseUrl + "jobs/" + jobId

	r, err := http.NewRequest("GET", jobUrl, nil)
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

	payload := &Steps{}
	derr := json.NewDecoder(res.Body).Decode(payload)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	var status string
	for _, v := range payload.Steps {
		// check if the strings match
		if v.Title == "retrieve-transactions" {
			status = v.Status
			break
		}
	}

	return status
}
