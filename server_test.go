package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestSetKeyRequest(t *testing.T) {
	client := &http.Client{}

	url := "http://0.0.0.0:5680/set"

	params := map[string]interface{}{
		"key":     "message",
		"value":   "greetings",
		"has_ttl": false,
		"exp":     -1,
	}

	res := makeJsonRequest(client, url, "POST", params)
	assert.NotEqual(t, res, nil)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestGetKeyRequest(t *testing.T) {
	client := &http.Client{}

	url := "http://0.0.0.0:5680/get"

	params := map[string]interface{}{
		"key": "message",
	}

	var kv KeyValue
	var bs []byte
	var err error
	//var n int
	response := makeJsonRequest(client, url, "GET", params)
	assert.NotEqual(t, response, nil)
	assert.Equal(t, response.StatusCode, http.StatusOK)

	defer response.Body.Close()

	bs, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	//t.Log(n)
	err = json.Unmarshal(bs, &kv)
	if err != nil {
		t.Error(err)
	}

	t.Log(kv.Value)

	assert.Equal(t, kv.Value, "greetings")
}

func TestKeysRequest(t *testing.T) {
	client := &http.Client{}

	url := "http://0.0.0.0:5680/keys"

	params := map[string]interface{}{
		"pattern": "mes[a-zA-Z]*e",
	}

	response := makeJsonRequest(client, url, "GET", params)
	assert.NotEqual(t, response, nil)
	assert.Equal(t, response.StatusCode, http.StatusOK)
}

func makeJsonRequest(client *http.Client, url string, method string, params map[string]interface{}) *http.Response {

	jsonStr, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
		return nil
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
		return nil
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}
