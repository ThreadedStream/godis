package main

import (
	"crypto/sha256"
	"regexp"
	"time"
)

type ValueStore struct {
	value       interface{}
	hashedField [32]byte
	ttl         int64
}

type Store struct {
	//Structure: Key -> [value, field, expiration time]
	dict map[string]ValueStore
}

func (s *Store) initStore() {
	s.dict = make(map[string]ValueStore)
}

//Should optimize it in future
//I guess that observables would be a good optimization decision. Not sure, though
func (s *Store) checkExpired() {
	for k, v := range s.dict {
		if v.ttl > time.Now().UnixNano() {
			delete(s.dict, k)
		}
	}
}

func (s *Store) SET(kv KeyValue) string {
	if kv.HasTtl {
		s.dict[kv.Key] = ValueStore{kv.Value, nil, -1}
	} else {
		s.dict[kv.Key] = ValueStore{kv.Value, nil, kv.Exp}
	}
	return "OK"
}

func (s *Store) HSET(kv KeyValue) string {
	hashedField := sha256.Sum256([]byte(kv.Field))
	if kv.HasTtl {
		s.dict[kv.Key] = ValueStore{kv.Value, hashedField, kv.Exp}
	} else {
		s.dict[kv.Key] = ValueStore{kv.Value, hashedField, -1}
	}

	return "OK"
}

func (s *Store) GET(key string) interface{} {
	value, ok := s.dict[key]
	if !ok {
		return nil
	}
	if value.ttl > time.Now().UnixNano() {
		delete(s.dict, key)
		return ""
	}

	return value.value
}

func (s *Store) DEL(key string) string {
	_, ok := s.dict[key]
	if ok {
		delete(s.dict, key)
		return "OK"
	}
	return ""
}

func (s *Store) KEYS(pattern string) []string {
	s.checkExpired()
	var list []string
	for k, _ := range s.dict {
		matched, _ := regexp.MatchString(pattern, k)
		if matched {
			list = append(list, k)
		}
	}
	return list
}
