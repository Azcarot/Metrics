package handlers

import (
	"bytes"
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

func AgentWorkers(data WorkerData) {
	sendAttempts := 3
	timeBeforeAttempt := 1
	err := PostJSONMetrics(data.BodyJSON, data.Batchrout, data.AgentflagData)
	for err != nil {
		if sendAttempts == 0 {
			break
		}
		times := time.Duration(timeBeforeAttempt)
		time.Sleep(times * time.Second)
		sendAttempts -= 1
		timeBeforeAttempt += 2

		PostJSONMetrics(data.BodyJSON, data.Batchrout, data.AgentflagData)

	}

	for _, buf := range data.Body {
		PostJSONMetrics(buf, data.Singlerout, data.AgentflagData)
	}
}

func PostJSONMetrics(b []byte, a string, f agentconfigs.AgentData) error {
	pth := "http://" + a
	var hashedMetrics string
	var err error
	b, err = agentconfigs.GzipForAgent(b)
	if err != nil {
		return err
	}
	resp, err := http.NewRequest("POST", pth, bytes.NewBuffer(b))
	if err != nil {
		return err
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
		defer res.Body.Close()
		return err
	}
	defer res.Body.Close()
	return err
}
