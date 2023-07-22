package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Transactions struct {
	Data  *[]TransactionData
	Links Links
}

type TransactionData struct {
	Status   string `json:"status"`
	Amount   string `json:"amount"`
	SubClass *SubClass
}

type SubClass struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type Links struct {
	Self string `json:"self"`
	Next string `json:"next"`
}

type AvgMap struct {
	title   string
	sum     float64
	count   int64
	average float64
}

func GetTransactionsList(accessToken string, userId string, params ...string) *[]TransactionData {
	nextQuery := ""
	if len(params) != 0 {
		nextQuery = "?" + params[0]
	}

	getTransactionsUrl := BaseUrl + "users/" + userId + "/transactions" + nextQuery

	r, err := http.NewRequest("GET", getTransactionsUrl, nil)
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

	payload := &Transactions{}
	derr := json.NewDecoder(res.Body).Decode(payload)
	if derr != nil {
		panic(derr)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
	var transactions *[]TransactionData = payload.Data
	if payload.Links.Next != "" {
		n := strings.Split(payload.Links.Next, "?")
		*transactions = append(*transactions, *GetTransactionsList(accessToken, userId, n[1])...)
	}
	return transactions
}

func CalculateAverage(transactions *[]TransactionData) {
	avgMap := map[string]*AvgMap{}

	for _, v := range *transactions {
		if v.SubClass != nil && v.Status == "posted" {
			fAmount, err := strconv.ParseFloat(v.Amount, 64)
			if err != nil {
				panic(err)
			}
			if mapKey, ok := avgMap[v.SubClass.Code]; ok {
				mapKey.count += 1
				mapKey.sum += fAmount
				mapKey.average =
					mapKey.sum / float64(mapKey.count)
			} else {
				new_val := AvgMap{v.SubClass.Title, fAmount, 1, fAmount}
				avgMap[v.SubClass.Code] = &new_val
			}
		}
	}

	for k, v := range avgMap {
		fmt.Println("Code: ", k, "Average:", v.average)
	}
}
