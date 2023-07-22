package main

import (
	"fmt"
	"time"

	api "basiq-task-go/api"
)

func main() {
	fmt.Println("Authorizing...")
	accessToken := api.AuthorizeUser()
	fmt.Println("Creating user...")
	userId := api.CreateUser(accessToken)
	fmt.Println("Connecting user...")
	connectJobId := api.ConnectUser(accessToken, userId)

	var jobStatus string
	timeoutCnt := 0
	for jobStatus != "success" {
		if timeoutCnt == 3 {
			panic("Job status timed out!")
		}
		if jobStatus == "failed" {
			panic("Connection failed!")
		}

		jobStatus = api.CheckConnectionJob(accessToken, connectJobId)
		time.Sleep(3 * time.Second)
		timeoutCnt += 1
	}
	fmt.Println("Connection successful! Fetching user transactions...")
	tList := api.GetTransactionsList(accessToken, userId)
	fmt.Println("Calculating average expenditures...")
	api.CalculateAverage(tList)
}
