package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func PostJSONMetrics(b []byte, a string, f agentconfigs.AgentData) (*http.Response, error) {
	pth := "http://" + a
	var hashedMetrics string

	if len(f.HashKey) > 0 {
		hashedMetrics = agentconfigs.MakeSHA(b, f.HashKey)
	}
	b = agentconfigs.GzipForAgent(b)

	resp, err := http.NewRequest("POST", pth, bytes.NewBuffer(b))
	if err != nil {
		panic(fmt.Sprintf("cannot post %s ", b))
	}
	if len(hashedMetrics) > 0 {
		resp.Header.Set("HashSHA256", hashedMetrics)
	}
	resp.Header.Add("Content-Type", storage.JSONContentType)
	resp.Header.Add("Content-Encoding", "gzip")
	client := &http.Client{}
	res, err := client.Do(resp)
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
	res, _ := client.Do(resp)
	return res
}
