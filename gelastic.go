package gelastic

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Gelastic struct {
	ElasticEndpoint *http.Client
	ElasticUser 	string
	ElasticPassword	string

	ElasticTimeFormat string

	Address			string
	Port			int
}


func (g *Gelastic) InitClient(Address string, Port int, User string, Pass string) {
	g.ElasticEndpoint = &http.Client{
		
	}

	g.Address 	= Address
	g.Port 		= Port
	g.ElasticUser = User
	g.ElasticPassword = Pass

	g.ElasticTimeFormat = "2006-01-02T15:04:05Z07:00"
}

func (g *Gelastic) makeElasticRequest(method string, path string, data []byte) {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/%s", g.Address, g.Port, path), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.ElasticUser, g.ElasticPassword)

	resp, err := g.ElasticEndpoint.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	str, _ := io.ReadAll(resp.Body)
	println(string(str))
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