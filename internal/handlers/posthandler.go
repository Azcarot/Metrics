package handlers

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"regexp"
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

// AgentWorkers отправляет собранные метрики на сервер,
// путь к серверу определяется соответствующим флагом в agentconfigs.AgentData
// осуществляет 3 попытки отправки, отправляет запрос раз в секунду
func AgentWorkers(data WorkerData) {
	sendAttempts := 3
	timeBeforeAttempt := 1
	var err error
	var encryptionKey []byte
	if data.AgentflagData.CryptoKey != "" {
		encryptionKey, err = agentconfigs.GetPublicKey(data.AgentflagData.CryptoKey)
		if err != nil {
			log.Fatal(err)
		}
		data.BodyJSON, err = agentconfigs.CypherData(encryptionKey, data.BodyJSON)
		if err != nil {
			log.Fatal(err)
		}
		for i, buf := range data.Body {
			data.Body[i], err = agentconfigs.CypherData(encryptionKey, buf)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = PostJSONMetrics(data.BodyJSON, data.Batchrout, data.AgentflagData)
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

// PostJSONMetrics формирует и отправляет запрос с полученной метрикой на сервер
// Адрес сервера определяется строкой a
// при наличии флага HashKey кодирует отправляемую метрику в sha256
func PostJSONMetrics(b []byte, a string, f agentconfigs.AgentData) error {
	pth := "http://" + a
	var hashedMetrics string
	var err error
	ip := GetOutboundIP(a)
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

	if f.CryptoKey != "" {
		resp.Header.Add("Crypto", "enabled")
	}
	resp.Header.Add("Content-Type", storage.JSONContentType)
	resp.Header.Add("Content-Encoding", "gzip")
	resp.Header.Add("X-Real-IP", ip.String())
	client := &http.Client{}
	res, err := client.Do(resp)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return err
}

func GetOutboundIP(a string) net.IP {
	regex := regexp.MustCompile("/")
	split := regex.Split(a, -1)
	conn, err := net.Dial("udp", split[0])
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
