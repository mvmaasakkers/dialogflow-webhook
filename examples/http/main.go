package main

import (
	"encoding/json"
	"log"
	"net/http"

	df "github.com/mvmaasakkers/dialogflow-webhook"
)

type params struct {
	City   string `json:"city"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func webhook(rw http.ResponseWriter, req *http.Request) {
	var err error
	var dfr *df.Request
	var p params

	decoder := json.NewDecoder(req.Body)
	if err = decoder.Decode(&dfr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Filter on action, using a switch for example

	// Retrieve the params of the request
	if err = dfr.GetParams(&p); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve a specific context
	if err = dfr.GetContext("my-awesome-context", &p); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Do things with the context you just retrieved
	dff := &df.Fulfillment{
		FulfillmentMessages: df.Messages{
			df.ForGoogle(df.SingleSimpleResponse("hello", "hello")),
			{RichMessage: df.Text{Text: []string{"hello"}}},
		},
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(dff)
}

func main() {
	http.HandleFunc("/webhook", webhook)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
