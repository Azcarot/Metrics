package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func PostJSONMetrics(b []byte, a string) *http.Response {
	pth := "http://" + a
	b = agentconfigs.GzipForAgent(b)
	resp, err := http.NewRequest("POST", pth, bytes.NewBuffer(b))
	if err != nil {
		panic(fmt.Sprintf("cannot post %s ", b))
	}
	resp.Header.Add("Content-Type", storage.JSONContentType)
	resp.Header.Add("Content-Encoding", "gzip")
	client := &http.Client{}
	res, _ := client.Do(resp)
	return res
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
