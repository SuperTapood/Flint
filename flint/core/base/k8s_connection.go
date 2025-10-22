package base

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type K8SConnection struct {
	Api   string
	Token string
}

func (connection *K8SConnection) Deploy(obj map[string]any) {
	data, _ := json.Marshal(obj)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var location = fmt.Sprint(obj["location"])
	req, err := http.NewRequest("POST", connection.Api+location, bytes.NewReader(data))
	req.Header.Add("Authorization", "Bearer "+connection.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))
}
