package main

import (
	"crypto/sha256"
	"encoding/json"
	"regexp"
	"time"
)

type ValueStore struct {
	Value       interface{} `json:"value"`
	IsHashed    bool        `json:"IsHashed"`
	HashedField [32]byte    `json:"HashedField"`
	Ttl         int64       `json:"ttl"`
}

type Store struct {
	//Structure: Key -> [value, field, expiration time]
	Dict map[string]ValueStore `json:"store"`
}

func (s *Store) initStore() {
	s.Dict = make(map[string]ValueStore)
}

func (s *Store) serializeStore() interface{} {
	bs, err := json.Marshal(s.Dict)
	if err != nil {
		return err.Error()
	}

	return bs
}

//Should optimize it in future
//I guess that observables would be a good optimization decision. Not sure, though
func (s *Store) checkExpired() {
	for k, v := range s.Dict {
		if time.Now().Unix() > v.Ttl && v.Ttl != -1 {
			delete(s.Dict, k)
		}
	}
}

func (s *Store) SET(kv KeyValue) string {
	if kv.HasTtl {
		s.Dict[kv.Key] = ValueStore{Value: kv.Value, Ttl: time.Now().Unix() + kv.Exp}
	} else {
		s.Dict[kv.Key] = ValueStore{Value: kv.Value, Ttl: -1}
	}
	return "OK"
}

func (s *Store) HSET(kv KeyValue) int {
	hashedField := sha256.Sum256([]byte(kv.Field))
	if kv.HasTtl {
		s.Dict[kv.Key] = ValueStore{kv.Value, true, hashedField, kv.Exp}
	} else {
		s.Dict[kv.Key] = ValueStore{kv.Value, true, hashedField, -1}
	}

	return 1
}

func (s *Store) HGET(key, field string) interface{} {
	hashedField := sha256.Sum256([]byte(field))
	value, ok := s.Dict[key]
	if ok {
		if time.Now().Unix() > value.Ttl && value.Ttl != -1 {
			delete(s.Dict, key)
			return ""
		}
		if value.HashedField == hashedField {
			return value.Value
		}
	}
	return nil
}

func (s *Store) GET(key string) interface{} {
	value, ok := s.Dict[key]
	if !ok {
		return nil
	}
	if time.Now().Unix() > value.Ttl && value.Ttl != -1 {
		delete(s.Dict, key)
		return ""
	}

	if value.IsHashed {
		return ""
	}

	return value.Value
}

func (s *Store) DEL(key string) string {
	_, ok := s.Dict[key]
	if ok {
		delete(s.Dict, key)
		return "OK"
	}
	return ""
}

func (s *Store) KEYS(pattern string) []string {
	s.checkExpired()
	var list []string
	for k, _ := range s.Dict {
		matched, _ := regexp.MatchString(pattern, k)
		if matched {
			list = append(list, k)
		}
	}
	return list
}

func (s *Store) MSET(kvs []KeyValue) string {
	for i := range kvs {
		s.SET(kvs[i])
	}

	return "OK"
}

func (s *Store) MGET(keys []string) []interface{} {
	var values = make([]interface{}, 0)
	for i := 0; i < len(keys); i++ {
		values = append(values, s.GET(keys[i]))
	}

	return values
}
