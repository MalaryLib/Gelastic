package Gelastic

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

type Gelastic struct {
	ElasticEndpoint *http.Client
	ElasticUser 	string
	ElasticPassword	string

	Address			string
	Port			int
}


func (g *Gelastic) InitClient(Address string, Port int) {
	g.ElasticEndpoint = &http.Client{
		
	}

	g.Address 	= Address
	g.Port 		= Port
}

func (g *Gelastic) makeElasticRequest(method string, path string, data []byte) {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/%s", g.Address, g.Port, path), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.ElasticUser, g.ElasticPassword)
}

func (g *Gelastic) AddDocument(index string, data any) {
	// first step is to marshal the struct 
	json, err := json.Marshal(data)
	if err != nil {
		// on error we warn the user, but do not panic so as to 
		// avoid losing logs.
		println("[Error] Document marshalling was not succesful while creating document.")
		return
	}

	h := sha256.New()
	h.Write(json)

	g.makeElasticRequest("PUT", fmt.Sprintf("%s/_doc/%x", index, h.Sum(nil)), json)
}