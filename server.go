package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	w               http.ResponseWriter
	statusCode      int
	responseMessage interface{}
}

func (r *Response) jsonResponse() {
	response, _ := json.Marshal(r.responseMessage)

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	r.w.Write(response)
}

func responseShortcut(w http.ResponseWriter, statusCode int, responseMessage interface{}) {
	response := Response{
		w:               w,
		statusCode:      statusCode,
		responseMessage: responseMessage,
	}

	response.jsonResponse()
}

func (s *Store) SetKey(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue

	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value := s.SET(kv)
	responseShortcut(w, http.StatusOK, map[string]interface{}{"status": value})
	return
}

func (s *Store) GetKey(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue
	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value := s.GET(kv.Key)
	responseShortcut(w, http.StatusOK, map[string]interface{}{"value": value})
}

func (s *Store) Keys(w http.ResponseWriter, r *http.Request) {
	var krm KeysRequestModel

	err := json.NewDecoder(r.Body).Decode(&krm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keys := s.KEYS(krm.Pattern)

	responseShortcut(w, http.StatusOK, map[string]interface{}{"keys": keys})
}

func (s *Store) Del(w http.ResponseWriter, r *http.Request) {
	var drm DelRequestModel

	err := json.NewDecoder(r.Body).Decode(&drm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status := s.DEL(drm.Key)
	responseShortcut(w, http.StatusOK, map[string]interface{}{"status": status})
}

func (s *Store) HSet(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue

	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := s.HSET(kv)
	responseShortcut(w, http.StatusOK, map[string]interface{}{"status": status})
}
