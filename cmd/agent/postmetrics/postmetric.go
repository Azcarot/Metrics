package postmetrics

import (
	"agent/measure"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

func Makepath(m measure.MemStorage) []string {
	var path []string
	pathscount := 0
	for name, value := range m.Gaugemem {

		path[pathscount] = "http://localhost:8080/update"
		path[pathscount] += "/gauge/" + name + "/" + strconv.FormatFloat(float64(value), 'g', -1, 64)
		pathscount++
	}
	for name, value := range m.Countermem {
		path[pathscount] = "http://localhost:8080/update"
		path[pathscount] = "/counter/" + name + "/" + strconv.Itoa(int(value))
		pathscount++
	}
	return path
}

func PostMetrics(pth string) *http.Response {
	data := []byte(pth)
	r := bytes.NewReader(data)
	resp, err := http.Post(pth, "text/plain", r)
	if err != nil {
		panic(fmt.Sprintf("cannot post %s ", data))
	}
	return resp
}
