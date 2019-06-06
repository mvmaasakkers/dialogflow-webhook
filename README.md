# Dialogflow Webhook

Simple library to create compatible DialogFlow v2 webhooks using Go.

This package is only intended to create webhooks, it doesn't implement the whole 
DialogFlow API.

This was originally forked from [github.com/leboncoin/dialogflow-go-webhook]()

# Go SDK

This project is used for structs and contexts for use in webhooks. There is a new [Go SDK](https://github.com/GoogleCloudPlatform/google-cloud-go/tree/master/dialogflow/apiv2) for Dialogflow's v2 API. See [this article](https://medium.com/leboncoin-engineering-blog/dialogflow-webhook-golang-and-protobuf-6269269f17f6) for a tutorial on how to use the protobuf definition to handle webhook.

# Table of Content

<!-- TOC -->

- [Dialogflow Webhook](#dialogflow-webhook)
- [Table of Content](#table-of-content)
- [Introduction](#introduction)
    - [Goal of this package](#goal-of-this-package)
    - [Disclaimer](#disclaimer)
- [Installation](#installation)
    - [Using dep](#using-dep)
    - [Using go get](#using-go-get)
- [Usage](#usage)
    - [Handling incoming request](#handling-incoming-request)
    - [Retrieving params and contexts](#retrieving-params-and-contexts)
    - [Responding with a fulfillment](#responding-with-a-fulfillment)
- [Examples](#examples)

<!-- /TOC -->

# Introduction

## Goal of this package

This package aims to implement a complete way to receive a DialogFlow payload,
parse it, and retrieve data stored in the parameters and contexts by providing
your own data structures. 

It also allows you to format your response properly, the way DialogFlow expects
it, including all the message types and platforms. (Such as cards, carousels,
quick replies, etc…)

## Disclaimer

As DialogFlow's v2 API is still in Beta, there may be breaking changes. If 
something breaks, please file an issue.

# Installation

## Using go get

If your project isn't using `dep` yet, you can use `go get` to install this
package :

`go get github.com/leboncoin/dialogflow-go-webhook`

# Usage

Import the package `github.com/mvmaasakkers/dialogflow-webhook` and use it with
the `dialogflow` package name. To make your code cleaner, import it as `df`.
All the following examples and usages use the `df` notation.

## Handling incoming request

When DialogFlow sends a request to your webhook, you can unmarshal the incoming
data to a `df.Request`. This, however, will not unmarshal the contexts and the
parameters because those are completely dependent on your data models.

```go
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

```

## Retrieving params and contexts

```go
type params struct {
	City   string `json:"city"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func webhook(c *gin.Context) {
	var err error
	var dfr *df.Request
	var p params

	if err = c.BindJSON(&dfr); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Retrieve the params of the request
	if err = dfr.GetParams(&p); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
```

In this example we're getting the DialogFlow request and unmarshalling the
params to a defined struct. This is why `json.RawMessage` is used in both the 
`Request.QueryResult.Parameters` and in the `Request.QueryResult.Contexts`.

This also allows you to filter and route according to the `action` and `intent`
DialogFlow detected, which means that depending on which action you detected,
you can unmarshal the parameters and contexts to a completely different data
structure.

The same thing can be done for contexts :

```go
type params struct {
	City   string `json:"city"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func webhook(c *gin.Context) {
	var err error
	var dfr *df.Request
	var p params

	if err = c.BindJSON(&dfr); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = dfr.GetContext("my-awesome-context", &p); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
```

## Responding with a fulfillment

DialogFlow expects you to respond with what is called a [fulfillment](https://dialogflow.com/docs/reference/api-v2/rest/v2beta1/WebhookResponse).

This package supports every rich response type.

```go
func webhook(c *gin.Context) {
	var err error
	var dfr *df.Request
	var p params

	if err = c.BindJSON(&dfr); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Send back a fulfillment
	dff := &df.Fulfillment{
		FulfillmentMessages: df.Messages{
			df.ForGoogle(df.SingleSimpleResponse("hello", "hello")),
			{RichMessage: df.Text{Text: []string{"hello"}}},
		},
	}
	c.JSON(http.StatusOK, dff)
}
```

