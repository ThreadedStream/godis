package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Parser struct {
	client *http.Client
}

func (p *Parser) parseCommand(command string, args []string) interface{} {
	var result interface{}

	switch command {
	case "SET":
		result = p.parseSET(args)
	case "GET":
		result = p.parseGET(args[0])
	case "KEYS":
		result = p.parseKEYS(args[0])
	case "MSET":
		result = p.parseMSET(args)
	case "HSET":
		result = p.parseHSET(args)
	case "HGET":
		result = p.parseHGET(args)
	case "DEL":
		result = p.parseDEL(args[0])
	default:
		result = "Unknown operation\n"
	}
	return result
}

func makeJsonRequest(client *http.Client, url string, method string, params interface{}) *http.Response {

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

func (p *Parser) initParserPipe() {
	p.client = &http.Client{}
}

func (p *Parser) parseSET(args []string) string {
	hasTtl := len(args) == 3
	var exp int
	var err error
	if !hasTtl {
		if len(args)%2 != 0 {
			return "Length of args should be divisible by 2"
		}
	} else {
		exp, err = strconv.Atoi(args[2])
		if err != nil {
			log.Println(err)
			return ""
		}
	}

	params := map[string]interface{}{
		"key":     args[0],
		"value":   args[1],
		"has_ttl": hasTtl,
		"exp":     exp,
	}

	url := "http://0.0.0.0:5680/set"
	method := "POST"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()

	var responseEntities ResponseEntities
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		log.Println(err)
		return ""
	}

	return responseEntities.Status
}

func (p *Parser) parseMSET(args []string) string {
	if len(args)%2 != 0 {
		return "Length of args should be divisible 2\n"
	}
	params := make([]map[string]interface{}, 0)

	//key value key1 value1 key2 value2

	index := 0
	for i := 0; i < len(args)/2; i++ {
		param := map[string]interface{}{
			"key":     args[index],
			"value":   args[index+1],
			"has_ttl": false,
			"exp":     -1,
		}
		params = append(params, param)
		index += 2
	}

	url := "http://0.0.0.0:5680/mset"

	method := "POST"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()

	var responseEntities ResponseEntities
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}
	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}

	return responseEntities.Status
}

func (p *Parser) parseHSET(args []string) interface{} {
	if len(args)%3 != 0 {
		return "Length of args should be divisible 3\n"
	}

	params := map[string]interface{}{
		"key":   args[0],
		"field": args[1],
		"value": args[2],
	}

	url := "http://0.0.0.0:5680/hset"
	method := "POST"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()

	var responseEntities ResponseEntities
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}
	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}

	return responseEntities.Created
}

func (p *Parser) parseGET(arg string) interface{} {

	params := map[string]interface{}{
		"key": arg,
	}

	url := "http://0.0.0.0:5680/get"
	method := "GET"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()
	var responseEntities ResponseEntities

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}

	return responseEntities.Value
}

func (p *Parser) parseHGET(args []string) interface{} {
	if len(args)%2 != 0 {
		return "Length of args should be divisible 2"
	}
	params := map[string]interface{}{
		"key":   args[0],
		"field": args[1],
	}

	url := "http://0.0.0.0:5680/hget"
	method := "GET"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()
	var responseEntities ResponseEntities

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}

	return responseEntities.Value
}

func (p *Parser) parseDEL(arg string) interface{} {
	params := map[string]interface{}{
		"key": arg,
	}
	url := "http://0.0.0.0:5680/del"
	method := "DELETE"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()
	var responseEntities ResponseEntities

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}

	return responseEntities.Status

}

//Arg signifies pattern
func (p *Parser) parseKEYS(arg string) interface{} {
	params := map[string]interface{}{
		"pattern": arg,
	}

	url := "http://0.0.0.0:5680/keys"

	method := "GET"

	response := makeJsonRequest(p.client, url, method, params)
	defer response.Body.Close()
	var responseEntities ResponseEntities

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	err = json.Unmarshal(bs, &responseEntities)
	if err != nil {
		return err.Error()
	}
	return p.prettifyKeysResponse(responseEntities.Keys)
}

func (p *Parser) prettifyKeysResponse(keys []string) string {
	var prettified string

	for i := range keys {
		prettified += fmt.Sprintf("%d) %s\n", i+1, keys[i])
	}

	return prettified
}
