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

func (connection *K8SConnection) Deploy(obj map[string]any, name string) {
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

	_, er := io.ReadAll(resp.Body)
	if er != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	//log.Println(string([]byte(body)))

	if resp.StatusCode == 409 {
		req, err := http.NewRequest("PUT", connection.Api+location+"/"+obj["metadata"].(map[string]any)["name"].(string), bytes.NewReader(data))
		req.Header.Add("Authorization", "Bearer "+connection.Token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()

		_, er := io.ReadAll(resp.Body)
		if er != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		//log.Println(string([]byte(body)))
	}
}
