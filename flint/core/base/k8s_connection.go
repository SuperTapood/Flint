package base

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	// Add this for pointer.Bool
)

type K8SConnection struct {
	Api   string
	Token string
}

func (connection *K8SConnection) makeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest(method, connection.Api+location, reader)

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IllEWm9BVHpfNTFDdkE0QkVWcjM2X21iVVRYYkhRZFY0bUVIUklCSzVRX1kifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYxNTQxNDU5LCJpYXQiOjE3NjE1Mzc4NTksImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiYmY5MzExNmUtY2NhMC00ZTg4LTk3YjUtMzMyN2EwNDg1MjQ3Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiMWJlNGM0NGItY2UyMy00MWRiLTg0ZWEtNzc5ZTU1NTUxYWZiIn19LCJuYmYiOjE3NjE1Mzc4NTksInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.DbUGoavPhUW3dnKNsLiD5q-3Q-JWdlx0MqXX6PWZwKJUDvAh24N3pOHx92prQW0gMvgiDiVPjE42RKGqOfI6PhLRyeGVkEZ7tgYhd9-ccn3gRIka4ogIrEBfsVP1TN-Uu5HzUg4cqSVniBPXxlcqxKw0UgbgHWbPE-SCYBSKQrLAQGU_qUSsIEQ83-Rw3HmypheaoSsb1YxRAVSGTMnvkiYWF054SGGc3fbzMMjmsfwFmodEdiBq5YJqzAJK5Rr9w0Pnb8XXS83aWE1O-W4e_Tg7cit8VV-7HBYehNj5ke8g8xlmoW8VCdEzjhRRf-scJIHhpgGbH4bayeHrB5BkSA")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, er := io.ReadAll(resp.Body)
	if er != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	fmt.Println(string(body))

	return body, resp
}

func mergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			if dm, ok := dv.(map[string]interface{}); ok {
				if sm, ok := v.(map[string]interface{}); ok {
					mergeMaps(dm, sm)
					continue
				}
			}
		}
		dst[k] = v
	}
}

func (connection *K8SConnection) Deploy(obj map[string]any, name string) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IllEWm9BVHpfNTFDdkE0QkVWcjM2X21iVVRYYkhRZFY0bUVIUklCSzVRX1kifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYxNzU5MzM2LCJpYXQiOjE3NjE3NTU3MzYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiMTI3NGEzNmItNjIwYS00ODhkLTkyZTMtMDRlZTc2MjM5OGM3Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiMWJlNGM0NGItY2UyMy00MWRiLTg0ZWEtNzc5ZTU1NTUxYWZiIn19LCJuYmYiOjE3NjE3NTU3MzYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.H8k3I5ZB4tLzT7vp9_CtO5jtj3kiu-0aX4Odll4ih2f3tk8IHlkHL21cnMLeWNaqNW3RGmreQzdq4thrARyWhctyEXM65tKLRyNfQtJy_i0GXwli8roUkv_ZG8zlEP1dcwZ9VXJyuVRn345yovT8Y8qdxRPyLnYx7DBSMZe0Y4Ko0bI1E4qRcn2yNtNcdiIVxEubWOLF1IsmkEPvCbz35tUIZHnMpri1w8gX6cf_DkXNpk4f_0QB_8zj9cd8h35i-zzJQbNMz7b9XbNqPnutvh1gJ1wytfURK4vomqkqhlITquY95lN7CfSHpbxNRrp6Ewopp31uh-C_encvvtmerA"

	delete(obj, "location")

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	config := &rest.Config{
		Host:        connection.Api,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	// g := &unstructured.Unstructured{
	// 	Object: obj,
	// }

	// unc := &unstructured.Unstructured{Object: obj}
	dynamicClient, _ := dynamic.NewForConfig(config)
	var g []byte
	g, er := json.Marshal(obj)
	fmt.Println(er)

	r, e := dynamicClient.Resource(gvr).Namespace("default").Patch(context.Background(), "nginx", types.MergePatchType, g, metav1.PatchOptions{})
	fmt.Println(r)
	fmt.Println(e)

	os.Exit(0)

	res, err := dynamicClient.Resource(gvr).Namespace("default").Get(context.Background(), "nginx", metav1.GetOptions{})
	fmt.Println(err)
	resultObject := res.Object
	delete(resultObject, "status")
	delete(resultObject["metadata"].(map[string]any), "managedFields")
	delete(resultObject["metadata"].(map[string]any), "creationTimestamp")
	delete(resultObject["metadata"].(map[string]any), "generation")
	delete(resultObject["metadata"].(map[string]any), "resourceVersion")
	delete(resultObject["metadata"].(map[string]any), "uid")
	resultObject["metadata"].(map[string]any)["labels"].(map[string]any)["nam"] = "ngn"
	fmt.Println(resultObject)
	// ob, err := res.MarshalJSON()
	// fmt.Println(err)
	// var o map[any]any
	// err = json.Unmarshal(ob, &o)
	// fmt.Println(err)
	// fmt.Println(o)
	// delete(o, "status")
	// fmt.Println(res)
	// fmt.Println(err)

	// OVERRIDE THE SPECIFIC FILEDS IN GET (TRYING TO MAINTAIN THE VOLUMES STUFF)

	data, _ := json.Marshal(resultObject)

	// dynamicClient.Resource(gvr).Namespace("default").Create(
	// 	context.Background(),
	// 	&unstructured.Unstructured{Object: obj},
	// 	metav1.CreateOptions{},
	// )

	// fmt.Println(a)
	// fmt.Println(b)

	c, d := dynamicClient.Resource(gvr).Namespace("default").Patch(
		context.Background(),
		"nginx",
		types.MergePatchType,
		data,
		metav1.PatchOptions{
			FieldManager: "my-client",
		},
	)

	fmt.Println(obj)
	fmt.Println(c)
	fmt.Println(d)

	os.Exit(0)
}

func (connection *K8SConnection) List() {
	var body, _ = connection.makeRequest("GET", "/api/v1/secrets", bytes.NewReader(make([]byte, 1)))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	secrets := make(map[string]any, 0)

	for i, secret := range result["items"].([]any) {
		fmt.Println(i)
		fmt.Println(secret.(map[string]any)["metadata"].(map[string]any)["name"])
		fmt.Println(secret.(map[string]any)["type"])
		if secret.(map[string]any)["type"] == "Opaque" {
			secrets[secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)] = secret
		}
	}
}
