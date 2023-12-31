package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Azcarot/Metrics/cmd/types"
)

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
