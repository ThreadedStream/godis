package main

import (
	"fmt"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	store := Store{}
	store.initStore()
	kv := KeyValue{
		Key:    "key",
		Value:  "some value",
		HasTtl: true,
		Exp:    120,
	}
	for i := 0; i < b.N; i++ {
		store.SET(kv)
	}
}

func BenchmarkHSet(b *testing.B) {
	store := Store{}
	store.initStore()
	kv := KeyValue{
		Key:    "key",
		Value:  "some value",
		Field:  "field",
		HasTtl: true,
		Exp:    120,
	}
	for i := 0; i < b.N; i++ {
		store.HSET(kv)
	}
}

func BenchmarkMSet(b *testing.B) {
	store := Store{}
	store.initStore()
	kvs := make([]KeyValue, 0)

	for i := 0; i < 3; i++ {
		kvs = append(kvs, KeyValue{
			Key:   fmt.Sprintf("Key%d", i),
			Value: fmt.Sprintf("Value%d", i),
		})
	}
	for i := 0; i < b.N; i++ {
		store.MSET(kvs)
	}
}

func BenchmarkGet(b *testing.B) {
	store := Store{}
	store.initStore()
	kv := KeyValue{
		Key: "key",
	}
	for i := 0; i < b.N; i++ {
		store.GET(kv.Key)
	}
}

func BenchmarkHGet(b *testing.B) {
	store := Store{}
	store.initStore()
	kv := KeyValue{
		Key:   "key",
		Field: "field",
	}
	for i := 0; i < b.N; i++ {
		store.HGET(kv.Key, kv.Field)
	}
}
func BenchmarkMGet(b *testing.B) {
	store := Store{}
	store.initStore()
	args := make([]string, 0)

	for i := 0; i < 3; i++ {
		args = append(args, fmt.Sprintf("key%d", i))
	}
	for i := 0; i < b.N; i++ {
		store.MGET(args)
	}
}

func BenchmarkKeys(b *testing.B) {
	store := Store{}
	store.initStore()
	pattern := "ke[yaz][1-9]*[1-9]*[a-zA-Z]+"

	for i := 0; i < b.N; i++ {
		store.KEYS(pattern)
	}
}

func BenchmarkDel(b *testing.B) {
	store := Store{}
	store.initStore()
	key := "key1"

	for i := 0; i < b.N; i++ {
		store.DEL(key)
	}
}

/*Conclusion:
SET: (num_of_ops: 18236211,speed: 62.6 ns/op)
GET: (num_of_ops: 136399186, speed: 8.67 ns/op)
HSET: (num_of_ops: 3156441, speed: 386 ns/op)
HGET: (num_of_ops: 3280952, speed:381 ns/op)
MSET: (num_of_ops: 20378941, speed: 58.5 ns/op)
MGET: (num_of_ops: 5733614, speed: 199 ns/op)
KEYS: (num_of_ops:100000000, speed:11.3 ns/op)
DEL: (num_of_ops:463427287, speed:2.64 ns/op)

Total time taken: 11.727s
*/
