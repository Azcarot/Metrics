package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Azcarot/Metrics/cmd/types"
)

func MakeJSON(m types.MemStorage) [][]byte {
	var body [][]byte
	var metric types.Metrics
	for name, value := range m.Gaugemem {
		newvalue := float64(value)
		metric.ID = name
		metric.MType = "gauge"
		metric.Value = &newvalue
		resp, err := json.Marshal(metric)
		if err != nil {
			panic(fmt.Sprintf("cannot make json %s ", body))
		}
		body = append(body, resp)
	}
	for name, value := range m.Countermem {
		newvalue := int64(value)
		metric.ID = name
		metric.MType = "counter"
		metric.Delta = &newvalue
		resp, err := json.Marshal(metric)
		if err != nil {
			panic(fmt.Sprintf("cannot make json %s ", body))
		}
		body = append(body, resp)

	}

	return body
}

func Makepath(m types.MemStorage, a string) []string {
	var path []string
	pathscount := 0
	for name, value := range m.Gaugemem {
		path = append(path, "http://"+a+"/update/gauge/"+name+"/"+strconv.FormatFloat(float64(value), 'g', -1, 64))
		pathscount++
	}
	for name, value := range m.Countermem {
		path = append(path, "http://"+a+"/update/counter/"+name+"/"+strconv.Itoa(int(value)))
		pathscount++
	}
	return path
}

func PostJSONMetrics(b []byte, a string) *http.Response {
	pth := "http://" + a + "/update/"
	resp, _ := http.NewRequest("POST", pth, bytes.NewBuffer(b))
	resp.Header.Add("Content-Type", types.CounterType)
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
