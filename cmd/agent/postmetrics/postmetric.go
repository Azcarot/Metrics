package postmetrics

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Azcarot/Metrics/cmd/agent/measure"
)

func Makepath(m measure.MemStorage) []string {
	var path []string
	pathscount := 0
	for name, value := range m.Gaugemem {
		path = append(path, "http://localhost:8080/update/gauge/"+name+"/"+strconv.FormatFloat(float64(value), 'g', -1, 64))
		pathscount++
	}
	for name, value := range m.Countermem {
		path = append(path, "http://localhost:8080/update/counter/"+name+"/"+strconv.Itoa(int(value)))
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
	resp.Header.Add("Content-Type", "text.plain")
	client := &http.Client{}
	res, _ := client.Do(resp)
	// defer res.Body.Close()
	return res
}
