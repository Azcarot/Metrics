package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

type WorkerData struct {
	Batchrout     string
	Singlerout    string
	Body          [][]byte
	BodyJSON      []byte
	AgentflagData agentconfigs.AgentData
}

func AgentWorkers(data WorkerData, results chan<- *http.Response) {
	sendAttempts := 3
	timeBeforeAttempt := 1
	resp, err := PostJSONMetrics(data.BodyJSON, data.Batchrout, data.AgentflagData)
	for err != nil {
		if sendAttempts == 0 {
			resp.Body.Close()
			panic(err)
		}

		times := time.Duration(timeBeforeAttempt)
		time.Sleep(times * time.Second)
		sendAttempts -= 1
		timeBeforeAttempt += 2
		resp.Body.Close()
		resp, err = PostJSONMetrics(data.BodyJSON, data.Batchrout, data.AgentflagData)
		if err != nil {
			resp.Body.Close()
			panic(err)
		}
		resp.Body.Close()

	}
	for _, buf := range data.Body {
		resp, _ = PostJSONMetrics(buf, data.Singlerout, data.AgentflagData)
		resp.Body.Close()
	}
	resp.Body.Close()
	results <- resp
	close(results)
}

func PostJSONMetrics(b []byte, a string, f agentconfigs.AgentData) (*http.Response, error) {
	pth := "http://" + a
	var hashedMetrics string

	b = agentconfigs.GzipForAgent(b)

	resp, err := http.NewRequest("POST", pth, bytes.NewBuffer(b))
	if err != nil {
		panic(fmt.Sprintf("cannot post %s ", b))
	}

	if len(f.HashKey) > 0 {
		hashedMetrics = agentconfigs.MakeSHA(b, f.HashKey)
		resp.Header.Add("HashSHA256", hashedMetrics)
	}
	resp.Header.Add("Content-Type", storage.JSONContentType)
	resp.Header.Add("Content-Encoding", "gzip")
	client := &http.Client{}
	res, err := client.Do(resp)
	if err != nil {
		panic("Cannot Post request")
	}
	defer res.Body.Close()
	return res, err
}

func PostMetrics(pth string) *http.Response {
	data := []byte(pth)
	resp, err := http.NewRequest("POST", pth, bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Sprintf("cannot post %s ", data))
	}
	resp.Header.Add("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(resp)
	if err != nil {
		panic("Cannot Post")
	}
	defer res.Body.Close()
	return res
}
